package cryptex

const (
	// System Secret
	CSize = 16

	// Salt
	SaltSize = 16

	// AES-GCM
	NonceSize  = 12
	TagSize    = 16
	AESKeySize = 32

	// Argon2id
	ArgonTime    = 3
	ArgonMemory  = 65536
	ArgonThreads = 4
	ArgonKeySize = 32

	// HKDF
	HKDFInfo = "vault-derive-km-v1"
)
