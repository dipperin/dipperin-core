package consts

const (
	// RecommendedSeedLen is the recommended length in bytes for a seed
	// to a master node.
	RecommendedSeedLen = 32 // 256 bits

	// HardenedKeyStart is the index at which a hardended key starts.  Each
	// extended key has 2^31 normal child keys and 2^31 hardned child keys.
	// Thus the range for normal child keys is [0, 2^31 - 1] and the range
	// for hardened child keys is [2^31, 2^32 - 1].
	HardenedKeyStart = 0x80000000 // 2^31

	// MinSeedBytes is the minimum number of bytes allowed for a seed to
	// a master node.
	MinSeedBytes = 16 // 128 bits

	// MaxSeedBytes is the maximum number of bytes allowed for a seed to
	// a master node.
	MaxSeedBytes = 64 // 512 bits

	// SerializedKeyLen is the length of a serialized public or private
	// extended key.  It consists of 4 bytes version, 1 byte depth, 4 bytes
	// fingerprint, 4 bytes child number, 32 bytes chain code, and 33 bytes
	// public/private key data.
	SerializedKeyLen = 4 + 1 + 4 + 4 + 32 + 33 // 78 bytes

	// MaxUint8 is the max positive integer which can be serialized in a uint8
	MaxUint8 = 1<<8 - 1
)

const ZerodExtendedKey  =  "zeroed extended key"

