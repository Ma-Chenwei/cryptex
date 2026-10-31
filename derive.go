package cryptex

import (
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
)

// deriveArgon2id derives a 32-byte key.
func deriveArgon2id(
	password []byte,
	salt []byte,
) []byte {
	return argon2.IDKey(
		password,
		salt,
		ArgonTime,
		ArgonMemory,
		ArgonThreads,
		ArgonKeySize,
	)
}

// deriveKM derives intermediate key material using HKDF-SHA256.
func deriveKM(
	password []byte,
	c []byte,
) ([]byte, error) {

	reader := hkdf.New(
		sha256.New,
		password,
		c,
		[]byte(HKDFInfo),
	)

	km := make([]byte, AESKeySize)

	if _, err := io.ReadFull(reader, km); err != nil {
		wipe(km)
		return nil, err
	}

	return km, nil
}
