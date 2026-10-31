package cryptex

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateC generates a random 16-byte system secret.
func GenerateC() ([]byte, error) {
	return randomBytes(CSize)
}

// randomBytes generates cryptographically secure random bytes.
func randomBytes(size int) ([]byte, error) {
	data := make([]byte, size)

	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		wipe(data)
		return nil, fmt.Errorf("%w: %v", ErrRandomGeneration, err)
	}

	return data, nil
}

// wipe attempts to clear sensitive byte slices.
func wipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
