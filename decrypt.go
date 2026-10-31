package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// Decrypt decrypts vault_cipher and restores the original maintext.
func Decrypt(
	m string,
	salt_a []byte,
	nonce_a []byte,
	s_cipher []byte,
	salt_b []byte,
	nonce_b []byte,
	vault_cipher []byte,
) (
	maintext []byte,
	err error,
) {

	// Validate master password.
	if m == "" {
		return nil, ErrInvalidPassword
	}

	// Validate salts.
	if len(salt_a) != SaltSize ||
		len(salt_b) != SaltSize {
		return nil, ErrInvalidSalt
	}

	// Validate nonces.
	if len(nonce_a) != NonceSize ||
		len(nonce_b) != NonceSize {
		return nil, ErrInvalidNonce
	}

	// Validate ciphertext lengths.
	if len(s_cipher) < TagSize ||
		len(vault_cipher) < TagSize {
		return nil, ErrInvalidCiphertext
	}

	password := []byte(m)

	// Derive Key A.
	keyA := deriveArgon2id(password, salt_a)
	defer wipe(keyA)

	blockA, err := aes.NewCipher(keyA)
	if err != nil {
		return nil, fmt.Errorf("cryptex: create AES cipher A: %w", err)
	}

	gcmA, err := cipher.NewGCM(blockA)
	if err != nil {
		return nil, fmt.Errorf("cryptex: create GCM A: %w", err)
	}

	// Recover C.
	c, err := gcmA.Open(
		nil,
		nonce_a,
		s_cipher,
		nil,
	)

	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	defer wipe(c)

	if len(c) != CSize {
		return nil, ErrInvalidCiphertext
	}

	// Derive KM.
	km, err := deriveKM(password, c)
	if err != nil {
		return nil, fmt.Errorf("cryptex: derive KM: %w", err)
	}

	defer wipe(km)

	// Derive Key B.
	keyB := deriveArgon2id(km, salt_b)
	defer wipe(keyB)

	blockB, err := aes.NewCipher(keyB)
	if err != nil {
		return nil, fmt.Errorf("cryptex: create AES cipher B: %w", err)
	}

	gcmB, err := cipher.NewGCM(blockB)
	if err != nil {
		return nil, fmt.Errorf("cryptex: create GCM B: %w", err)
	}

	// Recover maintext.
	maintext, err = gcmB.Open(
		nil,
		nonce_b,
		vault_cipher,
		nil,
	)

	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	return maintext, nil
}
