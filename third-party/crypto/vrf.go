package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/third-party/crypto/secp256k1"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
)

var (
	//curve = elliptic.P256()
	//params = curve.Params()
	curve = secp256k1.S256()

	// ErrInvalidVRF occurs when the VRF does not validate.
	// ErrEvalVRF occurs when unable to generate proof.
	ErrInvalidVRF = errors.New("invalid VRF proof")
	ErrEvalVRF    = errors.New("failed to evaluate vrf")
)

// hashToCurve hashes m to a curve point
func hashToCurve(m []byte) (x, y *big.Int) {
	h := sha256.New()
	var i uint32
	byteLen := (curve.BitSize + 7) >> 3
	for x == nil && i < 100 {
		// TODO: Use a NIST specified DRBG.
		h.Reset()
		binary.Write(h, binary.BigEndian, i)
		h.Write(m)
		r := []byte{2} // Set point encoding to "compressed", y=0.
		r = h.Sum(r)
		x, y = Unmarshal(curve, r[:byteLen+1])
		i++
	}
	return
}

var one = big.NewInt(1)

// hashToInt hashes to an integer [1,N-1]
func hashToInt(m []byte) *big.Int {
	// NIST SP 800-90A § A.5.1: Simple discard method.
	byteLen := (curve.BitSize + 7) >> 3
	h := sha256.New()
	for i := uint32(0); ; i++ {
		// TODO: Use a NIST specified DRBG.
		h.Reset()
		binary.Write(h, binary.BigEndian, i)
		h.Write(m)
		b := h.Sum(nil)
		k := new(big.Int).SetBytes(b[:byteLen])
		if k.Cmp(new(big.Int).Sub(curve.N, one)) == -1 {
			return k.Add(k, one)
		}
	}
}

// VRF returns the verifiable random function evaluated seed and an NIZK proof
// Check if you should use sha256.Sum256(vrf) of the output
func VRFProve(sk *ecdsa.PrivateKey, seed []byte) (vrf, nizk []byte, err error) {
	_, proof := Evaluate(sk, seed)
	if proof == nil {
		return nil, nil, ErrInvalidVRF
	}

	nizk = proof[0:64]
	vrf = proof[64 : 64+65]
	err = nil
	return
}

// Verify returns true if VRF and the NIZK proof is correct for seed
func VRFVerify(pk *ecdsa.PublicKey, seed, proof []byte) (bool, error) {
	_, err := ProofToHash(pk, seed, proof)
	if err != nil {
		return false, err
	}
	//arrayIdx, err := ProofToHash(pk, seed, proof)
	//if err != nil {
	//	return false, err
	//}
	//index := sha256.Sum256(vrf)
	//idx := arrayIdx[:]
	//idy := index[:]
	//fmt.Println(index)
	//fmt.Println(len(index))
	//if bytes.Compare(idx, idy) != 0 {
	//	return false, ErrInvalidVRF
	//}
	return true, nil
}

// Evaluate returns the verifiable unpredictable(random) function evaluated at seed
func Evaluate(sk *ecdsa.PrivateKey, seed []byte) (index [32]byte, proof []byte) {
	//r := sk.D.Bytes()
	//ri := sk.D

	nilIndex := [32]byte{}
	// Prover chooses r <-- [1,N-1]
	//r, _, _, err := elliptic.GenerateKey(curve, rand.Reader)
	// Same thing as ecdsa.GenerateKey
	r2, err := ecdsa.GenerateKey(curve, rand.Reader)
	r := r2.D.Bytes()
	if err != nil {
		return nilIndex, nil
	}
	ri := new(big.Int).SetBytes(r)

	// H = hashToCurve(m)
	Hx, Hy := hashToCurve(seed)

	// VRF_k(m) = [k]H
	sHx, sHy := curve.ScalarMult(Hx, Hy, sk.D.Bytes())
	vrf := elliptic.Marshal(curve, sHx, sHy) // 65 bytes.

	//// Test on curve
	//x := new(big.Int).SetBytes(vrf[1 : 1+32])
	//y := new(big.Int).SetBytes(vrf[1+32:])
	//fmt.Println(curve.IsOnCurve(x,y))

	// G is the base point
	// s = hashToInt(G, H, [k]G, VRF, [r]G, [r]H)
	rGx, rGy := curve.ScalarBaseMult(r)
	rHx, rHy := curve.ScalarMult(Hx, Hy, r)
	var b bytes.Buffer
	b.Write(elliptic.Marshal(curve, curve.Gx, curve.Gy))
	b.Write(elliptic.Marshal(curve, Hx, Hy))
	b.Write(elliptic.Marshal(curve, sk.PublicKey.X, sk.PublicKey.Y))
	b.Write(vrf)
	b.Write(elliptic.Marshal(curve, rGx, rGy))
	b.Write(elliptic.Marshal(curve, rHx, rHy))
	s := hashToInt(b.Bytes())

	// t = r−s*k mod N
	t := new(big.Int).Sub(ri, new(big.Int).Mul(s, sk.D))
	t.Mod(t, curve.N)

	// Index = H(vrf)
	index = sha256.Sum256(vrf)

	// Write s, t, and vrf to a proof blob. Also write leading zeros before s and t
	// if needed.
	var buf bytes.Buffer
	buf.Write(make([]byte, 32-len(s.Bytes())))
	buf.Write(s.Bytes())
	buf.Write(make([]byte, 32-len(t.Bytes())))
	buf.Write(t.Bytes())
	buf.Write(vrf)

	return index, buf.Bytes()
}

// ProofToHash asserts that proof is correct for seedz and outputs index.
func ProofToHash(pk *ecdsa.PublicKey, seed []byte, proof []byte) (index [32]byte, err error) {
	nilIndex := [32]byte{}
	// verifier checks that s == hashToInt(m, [t]G + [s]([k]G), [t]hashToCurve(m) + [s]VRF_k(m))
	if got, want := len(proof), 64+65; got != want {
		log.Info("vrp len check failed", "got", got, "want", want)
		return nilIndex, ErrInvalidVRF
	}

	// Parse proof into s, t, and vrf.
	s := proof[0:32]
	t := proof[32:64]
	// proof includes 'vrf' AKA hash
	vrf := proof[64 : 64+65]

	//// test
	//byteLen := (curve.BitSize + 7) >> 3
	//if len(vrf) != 1+2*byteLen {
	//	fmt.Println("AAA1")
	//}
	//if vrf[0] != 4 { // uncompressed form
	//	fmt.Println("AAA2")
	//}
	//p := curve.P
	//x := new(big.Int).SetBytes(vrf[1 : 1+byteLen])
	//y := new(big.Int).SetBytes(vrf[1+byteLen:])
	//fmt.Println(x,y)
	//if x.Cmp(p) >= 0 || y.Cmp(p) >= 0 {
	//	fmt.Println("A")
	//}
	//if !curve.IsOnCurve(x, y) {
	//	fmt.Println("AAA3")
	//}

	uHx, uHy := elliptic.Unmarshal(curve, vrf)
	if uHx == nil {
		log.Info("elliptic.Unmarshal(curve, vrf) failed")
		return nilIndex, ErrInvalidVRF
	}

	// [t]G + [s]([k]G) = [t+ks]G
	tGx, tGy := curve.ScalarBaseMult(t)
	ksGx, ksGy := curve.ScalarMult(pk.X, pk.Y, s)
	tksGx, tksGy := curve.Add(tGx, tGy, ksGx, ksGy)

	// H = hashToCurve(m)
	// [t]H + [s]VRF = [t+ks]H
	Hx, Hy := hashToCurve(seed)
	tHx, tHy := curve.ScalarMult(Hx, Hy, t)
	sHx, sHy := curve.ScalarMult(uHx, uHy, s)
	tksHx, tksHy := curve.Add(tHx, tHy, sHx, sHy)

	//   hashToInt(G, H, [k]G, VRF, [t]G + [s]([k]G), [t]H + [s]VRF)
	// = hashToInt(G, H, [k]G, VRF, [t+ks]G, [t+ks]H)
	// = hashToInt(G, H, [k]G, VRF, [r]G, [r]H)
	var b bytes.Buffer
	b.Write(elliptic.Marshal(curve, curve.Gx, curve.Gy))
	b.Write(elliptic.Marshal(curve, Hx, Hy))
	b.Write(elliptic.Marshal(curve, pk.X, pk.Y))
	b.Write(vrf)
	b.Write(elliptic.Marshal(curve, tksGx, tksGy))
	b.Write(elliptic.Marshal(curve, tksHx, tksHy))
	h2 := hashToInt(b.Bytes())

	// Left pad h2 with zeros if needed. This will ensure that h2 is padded
	// the same way s is.
	var buf bytes.Buffer
	buf.Write(make([]byte, 32-len(h2.Bytes())))
	buf.Write(h2.Bytes())

	if !hmac.Equal(s, buf.Bytes()) {
		log.Info("!hmac.Equal(s, buf.Bytes()) check failed", "s", hexutil.Encode(s), "buf.Bytes()", hexutil.Encode(buf.Bytes()))
		return nilIndex, ErrInvalidVRF
	}
	return sha256.Sum256(vrf), nil
}

//func VRFVerify(pk *ecdsa.PublicKey, seed []byte, index [32]byte, proof []byte) error {
//	arrayIdx, err := ProofToHash(pk, seed, proof)
//	if err != nil {
//		return err
//	}
//	idx := arrayIdx[:]
//	idy := index[:]
//	if bytes.Compare(idx, idy) != 0 {
//		return errors.New("error: index did not match")
//	}
//
//	return nil
//}
