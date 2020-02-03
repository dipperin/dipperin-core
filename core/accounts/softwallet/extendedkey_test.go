// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package softwallet

// References:
//   [BIP32]: BIP0032 - Hierarchical Deterministic Wallets
//   https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki

import (
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/stretchr/testify/assert"
	"math"
	"reflect"
	"testing"
)

type givenStruct struct {
	master string
	path   []uint32
}
type result struct {
	wantPub  string
	wantPriv string
}

// TestBIP0032Vectors tests the vectors provided by [BIP32] to ensure the
// derivation works as intended.

// todo  left to Integration Test
func TestBIP0032Vectors(t *testing.T) {
	// The master seeds for each of the two test vectors in [BIP32].
	testVec1MasterHex := "000102030405060708090a0b0c0d0e0f"
	testVec2MasterHex := "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542"
	testVec3MasterHex := "4b381541583be4423346c643850da4b320e46a87ae3d2a4e6da11eba819cd4acba45d239319ac14f863b8d5ab5a0d0c64d2e8a1e7d1457df2e5a3c51c73235be"
	hkStart := uint32(0x80000000)

	tests := []struct {
		name   string
		given  givenStruct
		expect result
	}{
		// Test vector 1
		{
			name:   "test vector 1 chain m",
			given:  givenStruct{testVec1MasterHex, []uint32{}},
			expect: result{"xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8", "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"},
		},
		{
			name:   "test vector 1 chain m/0H",
			given:  givenStruct{testVec1MasterHex, []uint32{hkStart}},
			expect: result{"xpub68Gmy5EdvgibQVfPdqkBBCHxA5htiqg55crXYuXoQRKfDBFA1WEjWgP6LHhwBZeNK1VTsfTFUHCdrfp1bgwQ9xv5ski8PX9rL2dZXvgGDnw", "xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7"},
		},
		{
			name: "test vector 1 chain m/0H/1", given: givenStruct{testVec1MasterHex, []uint32{hkStart, 1}},
			expect: result{"xpub6ASuArnXKPbfEwhqN6e3mwBcDTgzisQN1wXN9BJcM47sSikHjJf3UFHKkNAWbWMiGj7Wf5uMash7SyYq527Hqck2AxYysAA7xmALppuCkwQ", "xprv9wTYmMFdV23N2TdNG573QoEsfRrWKQgWeibmLntzniatZvR9BmLnvSxqu53Kw1UmYPxLgboyZQaXwTCg8MSY3H2EU4pWcQDnRnrVA1xe8fs"},
		},
		{
			name:  "test vector 1 chain m/0H/1/2H",
			given: givenStruct{testVec1MasterHex, []uint32{hkStart, 1, hkStart + 2}},
			expect: result{wantPub: "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
				wantPriv: "xprv9z4pot5VBttmtdRTWfWQmoH1taj2axGVzFqSb8C9xaxKymcFzXBDptWmT7FwuEzG3ryjH4ktypQSAewRiNMjANTtpgP4mLTj34bhnZX7UiM",
			}},
		{
			name:  "test vector 1 chain m/0H/1/2H/2",
			given: givenStruct{testVec1MasterHex, []uint32{hkStart, 1, hkStart + 2, 2}},
			expect: result{wantPub: "xpub6FHa3pjLCk84BayeJxFW2SP4XRrFd1JYnxeLeU8EqN3vDfZmbqBqaGJAyiLjTAwm6ZLRQUMv1ZACTj37sR62cfN7fe5JnJ7dh8zL4fiyLHV",
				wantPriv: "xprvA2JDeKCSNNZky6uBCviVfJSKyQ1mDYahRjijr5idH2WwLsEd4Hsb2Tyh8RfQMuPh7f7RtyzTtdrbdqqsunu5Mm3wDvUAKRHSC34sJ7in334",
			}},
		{
			name:  "test vector 1 chain m/0H/1/2H/2/1000000000",
			given: givenStruct{testVec1MasterHex, []uint32{hkStart, 1, hkStart + 2, 2, 1000000000}},
			expect: result{wantPub: "xpub6H1LXWLaKsWFhvm6RVpEL9P4KfRZSW7abD2ttkWP3SSQvnyA8FSVqNTEcYFgJS2UaFcxupHiYkro49S8yGasTvXEYBVPamhGW6cFJodrTHy",
				wantPriv: "xprvA41z7zogVVwxVSgdKUHDy1SKmdb533PjDz7J6N6mV6uS3ze1ai8FHa8kmHScGpWmj4WggLyQjgPie1rFSruoUihUZREPSL39UNdE3BBDu76",
			}},

		// Test vector 2
		{
			name:  "test vector 2 chain m",
			given: givenStruct{testVec2MasterHex, []uint32{}},
			expect: result{wantPub: "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB",
				wantPriv: "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U",
			}},
		{
			name:  "test vector 2 chain m/0",
			given: givenStruct{master: testVec2MasterHex, path: []uint32{0}},
			expect: result{wantPub: "xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH",
				wantPriv: "xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt",
			}},
		{
			name:  "test vector 2 chain m/0/2147483647H",
			given: givenStruct{master: testVec2MasterHex, path: []uint32{0, hkStart + 2147483647}},
			expect: result{wantPub: "xpub6ASAVgeehLbnwdqV6UKMHVzgqAG8Gr6riv3Fxxpj8ksbH9ebxaEyBLZ85ySDhKiLDBrQSARLq1uNRts8RuJiHjaDMBU4Zn9h8LZNnBC5y4a",
				wantPriv: "xprv9wSp6B7kry3Vj9m1zSnLvN3xH8RdsPP1Mh7fAaR7aRLcQMKTR2vidYEeEg2mUCTAwCd6vnxVrcjfy2kRgVsFawNzmjuHc2YmYRmagcEPdU9",
			}},
		{
			name:  "test vector 2 chain m/0/2147483647H/1",
			given: givenStruct{master: testVec2MasterHex, path: []uint32{0, hkStart + 2147483647, 1}},
			expect: result{wantPub: "xpub6DF8uhdarytz3FWdA8TvFSvvAh8dP3283MY7p2V4SeE2wyWmG5mg5EwVvmdMVCQcoNJxGoWaU9DCWh89LojfZ537wTfunKau47EL2dhHKon",
				wantPriv: "xprv9zFnWC6h2cLgpmSA46vutJzBcfJ8yaJGg8cX1e5StJh45BBciYTRXSd25UEPVuesF9yog62tGAQtHjXajPPdbRCHuWS6T8XA2ECKADdw4Ef",
			}},
		{
			name: "test vector 2 chain m/0/2147483647H/1/2147483646H",
			given: givenStruct{master: testVec2MasterHex,
				path: []uint32{0, hkStart + 2147483647, 1, hkStart + 2147483646},
			}, expect: result{wantPub: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
				wantPriv: "xprvA1RpRA33e1JQ7ifknakTFpgNXPmW2YvmhqLQYMmrj4xJXXWYpDPS3xz7iAxn8L39njGVyuoseXzU6rcxFLJ8HFsTjSyQbLYnMpCqE2VbFWc",
			}},
		{
			name: "test vector 2 chain m/0/2147483647H/1/2147483646H/2",
			given: givenStruct{master: testVec2MasterHex,
				path: []uint32{0, hkStart + 2147483647, 1, hkStart + 2147483646, 2},
			}, expect: result{wantPub: "xpub6FnCn6nSzZAw5Tw7cgR9bi15UV96gLZhjDstkXXxvCLsUXBGXPdSnLFbdpq8p9HmGsApME5hQTZ3emM2rnY5agb9rXpVGyy3bdW6EEgAtqt",
				wantPriv: "xprvA2nrNbFZABcdryreWet9Ea4LvTJcGsqrMzxHx98MMrotbir7yrKCEXw7nadnHM8Dq38EGfSh6dqA9QWTyefMLEcBYJUuekgW4BYPJcr9E7j",
			}},

		// Test vector 3
		{
			name: "test vector 3 chain m",
			given: givenStruct{master: testVec3MasterHex,
				path: []uint32{},
			}, expect: result{wantPub: "xpub661MyMwAqRbcEZVB4dScxMAdx6d4nFc9nvyvH3v4gJL378CSRZiYmhRoP7mBy6gSPSCYk6SzXPTf3ND1cZAceL7SfJ1Z3GC8vBgp2epUt13",
				wantPriv: "xprv9s21ZrQH143K25QhxbucbDDuQ4naNntJRi4KUfWT7xo4EKsHt2QJDu7KXp1A3u7Bi1j8ph3EGsZ9Xvz9dGuVrtHHs7pXeTzjuxBrCmmhgC6",
			}},
		{
			name: "test vector 3 chain m/0H",
			given: givenStruct{master: testVec3MasterHex,
				path: []uint32{hkStart},
			}, expect: result{wantPub: "xpub68NZiKmJWnxxS6aaHmn81bvJeTESw724CRDs6HbuccFQN9Ku14VQrADWgqbhhTHBaohPX4CjNLf9fq9MYo6oDaPPLPxSb7gwQN3ih19Zm4Y",
				wantPriv: "xprv9uPDJpEQgRQfDcW7BkF7eTya6RPxXeJCqCJGHuCJ4GiRVLzkTXBAJMu2qaMWPrS7AANYqdq6vcBcBUdJCVVFceUvJFjaPdGZ2y9WACViL4L",
			},
		}}

	for _, test := range tests {
		masterSeed, err := hex.DecodeString(test.given.master)
		assert.NoError(t, err)

		extKey, err := NewMaster(masterSeed, &DipperinChainCfg)
		assert.NoError(t, err)

		for _, childNum := range test.given.path {
			var err error
			extKey, err = extKey.Child(childNum)
			assert.NoError(t, err)
		}
		assert.Equal(t, extKey.Depth(), uint8(len(test.given.path)))

		privStr := extKey.String()

		pubKey, err := extKey.Neuter()
		assert.NoError(t, err)

		// Neutering a second time should have no effect.
		pubKey, err = pubKey.Neuter()
		assert.NoError(t, err)

		pubStr := pubKey.String()

		assert.Equal(t, test.expect.wantPriv, privStr)
		assert.Equal(t, test.expect.wantPub, pubStr)
	}
}

//TestPrivateDerivation tests several vectors which derive private keys from
//other private keys works as intended.
func TestPrivateDerivation(t *testing.T) {
	// The private extended keys for test vectors in [BIP32].
	testVec1MasterPrivKey := "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"
	testVec2MasterPrivKey := "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U"

	tests := []struct {
		name   string
		given  givenStruct
		expect result
	}{
		// Test vector 1
		{
			name: "test vector 1 chain m",
			given: givenStruct{master: testVec1MasterPrivKey,
				path: []uint32{}},
			expect: result{wantPriv: "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"}},
		{
			name: "test vector 1 chain m/0",
			given: givenStruct{master: testVec1MasterPrivKey,
				path: []uint32{0},
			}, expect: result{wantPriv: "xprv9uHRZZhbkedL37eZEnyrNsQPFZYRAvjy5rt6M1nbEkLSo378x1CQQLo2xxBvREwiK6kqf7GRNvsNEchwibzXaV6i5GcsgyjBeRguXhKsi4R"}},
		{
			name: "test vector 1 chain m/0/1",
			given: givenStruct{master: testVec1MasterPrivKey,
				path: []uint32{0, 1},
			}, expect: result{wantPriv: "xprv9ww7sMFLzJMzy7bV1qs7nGBxgKYrgcm3HcJvGb4yvNhT9vxXC7eX7WVULzCfxucFEn2TsVvJw25hH9d4mchywguGQCZvRgsiRaTY1HCqN8G"}},
		{
			name: "test vector 1 chain m/0/1/2",
			given: givenStruct{master: testVec1MasterPrivKey,
				path: []uint32{0, 1, 2},
			}, expect: result{wantPriv: "xprv9xrdP7iD2L1YZCgR9AecDgpDMZSTzP5KCfUykGXgjBxLgp1VFHsEeL3conzGAkbc1MigG1o8YqmfEA2jtkPdf4vwMaGJC2YSDbBTPAjfRUi"}},
		{
			name: "test vector 1 chain m/0/1/2/2",
			given: givenStruct{master: testVec1MasterPrivKey,
				path: []uint32{0, 1, 2, 2},
			}, expect: result{wantPriv: "xprvA2J8Hq4eiP7xCEBP7gzRJGJnd9CHTkEU6eTNMrZ6YR7H5boik8daFtDZxmJDfdMSKHwroCfAfsBKWWidRfBQjpegy6kzXSkQGGoMdWKz5Xh"}},
		{
			name: "test vector 1 chain m/0/1/2/2/1000000000",
			given: givenStruct{master: testVec1MasterPrivKey,
				path: []uint32{0, 1, 2, 2, 1000000000},
			},
			expect: result{wantPriv: "xprvA3XhazxncJqJsQcG85Gg61qwPQKiobAnWjuPpjKhExprZjfse6nErRwTMwGe6uGWXPSykZSTiYb2TXAm7Qhwj8KgRd2XaD21Styu6h6AwFz"}},

		// Test vector 2
		{
			name: "test vector 2 chain m",
			given: givenStruct{master: testVec2MasterPrivKey,
				path: []uint32{},
			},
			expect: result{wantPriv: "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U"}},
		{
			name: "test vector 2 chain m/0",
			given: givenStruct{master: testVec2MasterPrivKey,
				path: []uint32{0},
			},
			expect: result{wantPriv: "xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt"}},
		{
			name: "test vector 2 chain m/0/2147483647",
			given: givenStruct{master: testVec2MasterPrivKey,
				path: []uint32{0, 2147483647},
			},
			expect: result{wantPriv: "xprv9wSp6B7cXJWXZRpDbxkFg3ry2fuSyUfvboJ5Yi6YNw7i1bXmq9QwQ7EwMpeG4cK2pnMqEx1cLYD7cSGSCtruGSXC6ZSVDHugMsZgbuY62m6"}},
		{
			name: "test vector 2 chain m/0/2147483647/1",
			given: givenStruct{master: testVec2MasterPrivKey,
				path: []uint32{0, 2147483647, 1},
			},
			expect: result{wantPriv: "xprv9ysS5br6UbWCRCJcggvpUNMyhVWgD7NypY9gsVTMYmuRtZg8izyYC5Ey4T931WgWbfJwRDwfVFqV3b29gqHDbuEpGcbzf16pdomk54NXkSm"}},
		{
			name: "test vector 2 chain m/0/2147483647/1/2147483646",
			given: givenStruct{master: testVec2MasterPrivKey,
				path: []uint32{0, 2147483647, 1, 2147483646},
			},
			expect: result{wantPriv: "xprvA2LfeWWwRCxh4iqigcDMnUf2E3nVUFkntc93nmUYBtb9rpSPYWa8MY3x9ZHSLZkg4G84UefrDruVK3FhMLSJsGtBx883iddHNuH1LNpRrEp"}},
		{
			name: "test vector 2 chain m/0/2147483647/1/2147483646/2",
			given: givenStruct{master: testVec2MasterPrivKey,
				path: []uint32{0, 2147483647, 1, 2147483646, 2},
			},
			expect: result{wantPriv: "xprvA48ALo8BDjcRET68R5RsPzF3H7WeyYYtHcyUeLRGBPHXu6CJSGjwW7dWoeUWTEzT7LG3qk6Eg6x2ZoqD8gtyEFZecpAyvchksfLyg3Zbqam"}},

		// Custom tests to trigger specific conditions.
		{
			// Seed 000000000000000000000000000000da.
			name: "Derived privkey with zero high byte m/0",
			given: givenStruct{master: "xprv9s21ZrQH143K4FR6rNeqEK4EBhRgLjWLWhA3pw8iqgAKk82ypz58PXbrzU19opYcxw8JDJQF4id55PwTsN1Zv8Xt6SKvbr2KNU5y8jN8djz",
				path: []uint32{0},
			},
			expect: result{wantPriv: "xprv9uC5JqtViMmgcAMUxcsBCBFA7oYCNs4bozPbyvLfddjHou4rMiGEHipz94xNaPb1e4f18TRoPXfiXx4C3cDAcADqxCSRSSWLvMBRWPctSN9"}},
	}

	for _, test := range tests {
		extKey, err := NewKeyFromString(test.given.master)
		assert.NoError(t, err)

		for _, childNum := range test.given.path {
			var err error
			extKey, err = extKey.Child(childNum)
			assert.NoError(t, err)
		}

		privStr := extKey.String()

		assert.Equal(t, test.expect.wantPriv, privStr)
	}
}

// TestPublicDerivation tests several vectors which derive public keys from
// other public keys works as intended.
func TestPublicDerivation(t *testing.T) {
	// The public extended keys for test vectors in [BIP32].
	testVec1MasterPubKey := "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"
	testVec2MasterPubKey := "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB"

	tests := []struct {
		name   string
		given  givenStruct
		expect result
	}{
		// Test vector 1
		{
			name: "test vector 1 chain m",
			given: givenStruct{master: testVec1MasterPubKey,
				path: []uint32{}},
			expect: result{wantPub: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"}},
		{
			name: "test vector 1 chain m/0",
			given: givenStruct{master: testVec1MasterPubKey,
				path: []uint32{0},
			}, expect: result{wantPub: "xpub68Gmy5EVb2BdFbj2LpWrk1M7obNuaPTpT5oh9QCCo5sRfqSHVYWex97WpDZzszdzHzxXDAzPLVSwybe4uPYkSk4G3gnrPqqkV9RyNzAcNJ1"}},
		{
			name: "test vector 1 chain m/0/1",
			given: givenStruct{master: testVec1MasterPubKey,
				path: []uint32{0, 1},
			}, expect: result{wantPub: "xpub6AvUGrnEpfvJBbfx7sQ89Q8hEMPM65UteqEX4yUbUiES2jHfjexmfJoxCGSwFMZiPBaKQT1RiKWrKfuDV4vpgVs4Xn8PpPTR2i79rwHd4Zr"}},
		{
			name: "test vector 1 chain m/0/1/2",
			given: givenStruct{master: testVec1MasterPubKey,
				path: []uint32{0, 1, 2},
			}, expect: result{wantPub: "xpub6BqyndF6rhZqmgktFCBcapkwubGxPqoAZtQaYewJHXVKZcLdnqBVC8N6f6FSHWUghjuTLeubWyQWfJdk2G3tGgvgj3qngo4vLTnnSjAZckv"}},
		{
			name: "test vector 1 chain m/0/1/2/2",
			given: givenStruct{master: testVec1MasterPubKey,
				path: []uint32{0, 1, 2, 2},
			}, expect: result{wantPub: "xpub6FHUhLbYYkgFQiFrDiXRfQFXBB2msCxKTsNyAExi6keFxQ8sHfwpogY3p3s1ePSpUqLNYks5T6a3JqpCGszt4kxbyq7tUoFP5c8KWyiDtPp"}},
		{
			name: "test vector 1 chain m/0/1/2/2/1000000000",
			given: givenStruct{master: testVec1MasterPubKey,
				path: []uint32{0, 1, 2, 2, 1000000000},
			}, expect: result{wantPub: "xpub6GX3zWVgSgPc5tgjE6ogT9nfwSADD3tdsxpzd7jJoJMqSY12Be6VQEFwDCp6wAQoZsH2iq5nNocHEaVDxBcobPrkZCjYW3QUmoDYzMFBDu9"}},

		// Test vector 2
		{
			name: "test vector 2 chain m",
			given: givenStruct{master: testVec2MasterPubKey,
				path: []uint32{},
			}, expect: result{wantPub: "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB"}},
		{
			name: "test vector 2 chain m/0",
			given: givenStruct{master: testVec2MasterPubKey,
				path: []uint32{0},
			}, expect: result{wantPub: "xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH"}},
		{
			name: "test vector 2 chain m/0/2147483647",
			given: givenStruct{master: testVec2MasterPubKey,
				path: []uint32{0, 2147483647},
			}, expect: result{wantPub: "xpub6ASAVgeWMg4pmutghzHG3BohahjwNwPmy2DgM6W9wGegtPrvNgjBwuZRD7hSDFhYfunq8vDgwG4ah1gVzZysgp3UsKz7VNjCnSUJJ5T4fdD"}},
		{
			name: "test vector 2 chain m/0/2147483647/1",
			given: givenStruct{master: testVec2MasterPubKey,
				path: []uint32{0, 2147483647, 1},
			}, expect: result{wantPub: "xpub6CrnV7NzJy4VdgP5niTpqWJiFXMAca6qBm5Hfsry77SQmN1HGYHnjsZSujoHzdxf7ZNK5UVrmDXFPiEW2ecwHGWMFGUxPC9ARipss9rXd4b"}},
		{
			name: "test vector 2 chain m/0/2147483647/1/2147483646",
			given: givenStruct{master: testVec2MasterPubKey,
				path: []uint32{0, 2147483647, 1, 2147483646},
			}, expect: result{wantPub: "xpub6FL2423qFaWzHCvBndkN9cbkn5cysiUeFq4eb9t9kE88jcmY63tNuLNRzpHPdAM4dUpLhZ7aUm2cJ5zF7KYonf4jAPfRqTMTRBNkQL3Tfta"}},
		{
			name: "test vector 2 chain m/0/2147483647/1/2147483646/2",
			given: givenStruct{master: testVec2MasterPubKey,
				path: []uint32{0, 2147483647, 1, 2147483646, 2},
			}, expect: result{wantPub: "xpub6H7WkJf547AiSwAbX6xsm8Bmq9M9P1Gjequ5SipsjipWmtXSyp4C3uwzewedGEgAMsDy4jEvNTWtxLyqqHY9C12gaBmgUdk2CGmwachwnWK"}},
	}

	for _, test := range tests {
		extKey, err := NewKeyFromString(test.given.master)
		assert.NoError(t, err)

		for _, childNum := range test.given.path {
			var err error
			extKey, err = extKey.Child(childNum)
			assert.NoError(t, err)
		}

		pubStr := extKey.String()
		assert.Equal(t, test.expect.wantPub, pubStr)
	}
}

// TestGenenerateSeed ensures the GenerateSeed function works as intended.
func TestGenenerateSeed(t *testing.T) {
	tests := []struct {
		name   string
		given  uint8
		expect error
	}{
		// Test various valid lengths.
		{name: "test 16 bytes", given: 16},
		{name: "test 17 bytes", given: 17},
		{name: "test 20 bytes", given: 20},
		{name: "test 32 bytes", given: 32},
		{name: "test 64 bytes", given: 64},

		// Test invalid lengths.
		{name: "test 15 bytes", given: 15, expect: gerror.ErrInvalidSeedLen},
		{name: "test 65 bytes", given: 65, expect: gerror.ErrInvalidSeedLen},
	}

	for _, test := range tests {
		seed, err := GenerateSeed(test.given)
		//t.Log(err)
		//t.Log(test.expect)
		assert.True(t, reflect.DeepEqual(err, test.expect))
		if test.expect == nil {
			assert.Equal(t, int(test.given), len(seed))
		}
	}
}

// TestExtendedKeyAPI ensures the API on the ExtendedKey type works as intended.
func TestExtendedKeyAPI(t *testing.T) {

	type result struct {
		isPrivate  bool
		parentFP   uint32
		privKey    string
		privKeyErr error
		pubKey     string
	}

	tests := []struct {
		name   string
		given  string
		expect result
	}{
		{
			name:  "test vector 1 master node private",
			given: "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi",
			expect: result{isPrivate: true,
				parentFP: 0,
				privKey:  "e8f32e723decf4051aefac8e2c93c9c5b214313817cdb01a1494b917c8436b35",
				pubKey:   "0339a36013301597daef41fbe593a02cc513d0b55527ec2df1050e2e8ff49c85c2",
			},
		},
		{
			name:  "test vector 1 chain m/0H/1/2H public",
			given: "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
			expect: result{isPrivate: false,
				parentFP:   3203769081,
				privKeyErr: gerror.ErrNotPrivExtKey,
				pubKey:     "0357bfe1e341d01c69fe5654309956cbea516822fba8a601743a012a7896ee8dc2"},
		},
	}

	for _, test := range tests {
		key, err := NewKeyFromString(test.given)
		assert.NoError(t, err)
		assert.Equal(t, key.IsPrivate(), test.expect.isPrivate)

		parentFP := key.ParentFingerprint()
		assert.Equal(t, test.expect.parentFP, parentFP)

		serializedKey := key.String()
		assert.Equal(t, test.given, serializedKey)

		privKey, err := key.ECPrivKey()
		assert.True(t, reflect.DeepEqual(err, test.expect.privKeyErr))

		if test.expect.privKeyErr == nil {
			privKeyStr := hex.EncodeToString(privKey.Serialize())
			assert.Equal(t, test.expect.privKey, privKeyStr)
		}

		pubKey, err := key.ECPubKey()
		assert.NoError(t, err)

		pubKeyStr := hex.EncodeToString(pubKey.SerializeCompressed())
		assert.Equal(t, test.expect.pubKey, pubKeyStr)
	}
}

func Test_NewKeyFromString(t *testing.T) {

	type given struct {
		key    string
		neuter bool
	}

	type result struct {
		err       error
		neuterErr error
	}

	// NewKeyFromString failure tests.
	tests := []struct {
		name   string
		given  given
		expect result
	}{
		{
			name:   "invalid key expect",
			given:  given{key: "xpub1234"},
			expect: result{err: gerror.ErrInvalidKeyLen},
		},
		{
			name:   "bad checksum",
			given:  given{key: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EBygr15"},
			expect: result{err: gerror.ErrBadChecksum},
		},
		{
			name:   "pubkey not on curve",
			given:  given{key: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ1hr9Rwbk95YadvBkQXxzHBSngB8ndpW6QH7zhhsXZ2jHyZqPjk"},
			expect: result{err: errors.New("invalid square root")},
		},
		{
			name: "unsupported version",
			given: given{key: "xbad4LfUL9eKmA66w2GJdVMqhvDmYGJpTGjWRAtjHqoUY17sGaymoMV9Cm3ocn9Ud6Hh2vLFVC7KSKCRVVrqc6dsEdsTjRV1WUmkK85YEUujAPX",
				neuter: true},
			expect: result{err: nil,
				neuterErr: chaincfg.ErrUnknownHDKeyID},
		},
	}

	for _, test := range tests {
		extKey, err := NewKeyFromString(test.given.key)
		assert.True(t, reflect.DeepEqual(err, test.expect.err))

		if test.given.neuter {
			_, err := extKey.Neuter()
			assert.True(t, reflect.DeepEqual(err, test.expect.neuterErr))
		}
	}
}

func TestNewMaster(t *testing.T) {

	tests := []struct {
		name   string
		given  string
		expect error
	}{
		{
			name: "NewMasterRight",
			given:  "000102030405060708090a0b0c0d0e0f",
			expect: nil,
		},

		{
			name: "ErrInvalidSeedLen",
			given: "abfffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542",
			expect: gerror.ErrInvalidSeedLen,
		},
	}

	for _, test := range tests {
		t.Log(test.name)
		// Create new key from seed and get the neutered version.
		masterSeed, err := hex.DecodeString(test.given)
		assert.NoError(t, err)

		_, err = NewMaster(masterSeed, &DipperinChainCfg)
		if err != nil {
			assert.Equal(t, test.expect, err)
		}
	}
}

func TestExtendedKey_Child(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() (*ExtendedKey, uint32)
		expect error
	}{
		{
			name: "ErrDeriveBeyondMaxDepth",
			given: func() (*ExtendedKey, uint32) {
				extKey, err := NewMaster([]byte(`abcd1234abcd1234abcd1234abcd1234`), &DipperinChainCfg)
				assert.NoError(t, err)
				extKey.depth = math.MaxUint8
				return extKey,1
			},
			expect: gerror.ErrDeriveBeyondMaxDepth,
		},

		{
			name: "ErrDeriveHardFromPublic",
			given: func() (*ExtendedKey, uint32) {
				extKey, err := NewMaster([]byte(`abcd1234abcd1234abcd1234abcd1234`), &DipperinChainCfg)
				assert.NoError(t, err)
				extKey.isPrivate = false
				return extKey,consts.HardenedKeyStart
			},
			expect: gerror.ErrDeriveHardFromPublic,
		},
		{
			name: "ChildRight",
			given: func() (*ExtendedKey, uint32) {
				extKey, err := NewMaster([]byte(`abcd1234abcd1234abcd1234abcd1234`), &DipperinChainCfg)
				assert.NoError(t, err)
				return extKey,1
			},
			expect: nil,
		},
	}

	for _, tt := range testCases {
		t.Log(tt.name)
		key, i := tt.given()
		_, err := key.Child(i)

		if err != nil {
			assert.Equal(t, tt.expect, err)
		}
	}
}

func TestExtendedKey_Neuter(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() (*ExtendedKey)
		expect error
	}{
		{
			name: "ErrUnknownHDKeyID",
			given: func() (*ExtendedKey) {
				extKey, err := NewMaster([]byte(`abcd1234abcd1234abcd1234abcd1234`), &DipperinChainCfg)
				assert.NoError(t, err)
				extKey.version = []byte{12,23,34}
				return extKey
			},
			expect: chaincfg.ErrUnknownHDKeyID,
		},

		{
			name: "IsPublic",
			given: func() (*ExtendedKey) {
				extKey, err := NewMaster([]byte(`abcd1234abcd1234abcd1234abcd1234`), &DipperinChainCfg)
				assert.NoError(t, err)
				extKey.isPrivate = false
				return extKey
			},
			expect: nil,
		},
		{
			name: "NeuterRight",
			given: func() (*ExtendedKey) {
				extKey, err := NewMaster([]byte(`abcd1234abcd1234abcd1234abcd1234`), &DipperinChainCfg)
				assert.NoError(t, err)
				return extKey
			},
			expect: nil,
		},
	}

	for _, tt := range testCases {
		t.Log(tt.name)
		key := tt.given()
		neuterKey, err := key.Neuter()

		if err != nil {
			assert.Equal(t, tt.expect, err)
		} else if !key.IsPrivate() {
			assert.Equal(t, neuterKey, key)
		}
	}
}

// TestZeroExtendedKey ensures that zeroing an extended key works as intended.
func TestZeroExtendedKey(t *testing.T) {
	type given struct {
		master string
		extKey string
	}

	tests := []struct {
		name   string
		given  given
		expect bool
		//net    *chaincfg.Params
	}{
		// Test vector 1
		{
			name: "TestZeroRight",
			given: given{master: "000102030405060708090a0b0c0d0e0f",
				extKey: "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"},
			expect: true,
		},

		// Test vector 2
		{
			name: "TestZeroRightTwo",
			given: given{master: "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542",
				extKey: "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U"},
			expect: true,
		},
	}

	// Use a closure to test that a key is zeroed since the tests create
	// keys in different ways and need to test the same things multiple
	// times.
	testZeroed := func(testName string, key *ExtendedKey) bool {
		// Zeroing a key should result in it no longer being private
		assert.True(t, !key.IsPrivate())

		parentFP := key.ParentFingerprint()
		assert.Equal(t, uint32(0), parentFP)

		serializedKey := key.String()
		assert.Equal(t, consts.ZerodExtendedKey, serializedKey)

		_, err := key.ECPrivKey()
		assert.True(t, reflect.DeepEqual(err, gerror.ErrNotPrivExtKey))

		wantErr := errors.New("pubkey string is empty")
		_, err = key.ECPubKey()
		assert.Equal(t, wantErr, err)
		return true
	}

	for _, test := range tests {
		// Create new key from seed and get the neutered version.
		masterSeed, err := hex.DecodeString(test.given.master)
		assert.NoError(t, err)

		key, err := NewMaster(masterSeed, &DipperinChainCfg)
		assert.NoError(t, err)

		neuteredKey, err := key.Neuter()
		assert.NoError(t, err)

		// Ensure both non-neutered and neutered keys are zeroed
		// properly.
		key.Zero()
		assert.Equal(t, test.expect, testZeroed(test.name+" from seed not neutered", key))
		neuteredKey.Zero()
		assert.Equal(t, test.expect, testZeroed(test.name+" from seed neutered", key))
	}
}

// TestMaximumDepth ensures that attempting to retrieve a child key when already
// at the maximum depth is not allowed.  The serialization of a BIP32 key uses
// uint8 to encode the depth.  This implicitly bounds the depth of the tree to
// 255 derivations.  Here we test that an error is returned after 'max uint8'.
func TestMaximumDepth(t *testing.T) {
	extKey, err := NewMaster([]byte(`abcd1234abcd1234abcd1234abcd1234`), &DipperinChainCfg)
	assert.NoError(t, err)

	for i := uint8(0); i < math.MaxUint8; i++ {
		assert.Equal(t, extKey.Depth(), i)
		newKey, err := extKey.Child(1)
		assert.NoError(t, err)
		extKey = newKey
	}

	noKey, err := extKey.Child(1)
	assert.Equal(t, gerror.ErrDeriveBeyondMaxDepth, err)
	assert.Equal(t, (*ExtendedKey)(nil), noKey)

}

func TestNewExtendedKey(t *testing.T) {
	testSk := "xprv9wTYmMFdV23N2TdNG573QoEsfRrWKQgWeibmLntzniatZvR9BmLnvSxqu53Kw1UmYPxLgboyZQaXwTCg8MSY3H2EU4pWcQDnRnrVA1xe8fs"
	key, err := NewKeyFromString(testSk)

	assert.NoError(t, err)
	assert.Equal(t, uint8(2), key.Depth())
	assert.Equal(t, true, key.IsPrivate())
	assert.Equal(t, uint32(1545328200), key.ParentFingerprint())
}
