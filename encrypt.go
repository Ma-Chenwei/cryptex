package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// Encrypt encrypts maintext using the master password and system secret.
//
// If c is nil or empty, a new system secret is generated automatically.
//
// Returns:
// generatedC, salt_a, nonce_a, s_cipher,
// salt_b, nonce_b, vault_cipher, error
func Encrypt(
	m string,
	c []byte,
	maintext []byte,
) (
	generatedC []byte,
	salt_a []byte,
	nonce_a []byte,
	s_cipher []byte,
	salt_b []byte,
	nonce_b []byte,
	vault_cipher []byte,
	err error,
) {

	// Validate master password.
	if m == "" {
		err = ErrInvalidPassword
		return
	}

	// Select system secret.
	switch len(c) {
	case 0:
		generatedC, err = GenerateC()
		if err != nil {
			return
		}

	case CSize:
		generatedC = append([]byte(nil), c...)

	default:
		err = ErrInvalidKey
		return
	}

	// Generate salts.
	salt_a, err = randomBytes(SaltSize)
	if err != nil {
		return
	}

	salt_b, err = randomBytes(SaltSize)
	if err != nil {
		return
	}

	// Generate nonces.
	nonce_a, err = randomBytes(NonceSize)
	if err != nil {
		return
	}

	nonce_b, err = randomBytes(NonceSize)
	if err != nil {
		return
	}

	password := []byte(m)

	// Derive Key A.
	keyA := deriveArgon2id(password, salt_a)
	defer wipe(keyA)

	blockA, aesErr := aes.NewCipher(keyA)
	if aesErr != nil {
		err = fmt.Errorf("cryptex: create AES cipher A: %w", aesErr)
		return
	}

	gcmA, aesErr := cipher.NewGCM(blockA)
	if aesErr != nil {
		err = fmt.Errorf("cryptex: create GCM A: %w", aesErr)
		return
	}

	// Encrypt C.
	s_cipher = gcmA.Seal(
		nil,
		nonce_a,
		generatedC,
		nil,
	)

	// Derive intermediate key material KM.
	km, deriveErr := deriveKM(password, generatedC)
	if deriveErr != nil {
		err = fmt.Errorf("cryptex: derive KM: %w", deriveErr)
		return
	}
	defer wipe(km)

	// Derive Key B.
	keyB := deriveArgon2id(km, salt_b)
	defer wipe(keyB)

	blockB, aesErr := aes.NewCipher(keyB)
	if aesErr != nil {
		err = fmt.Errorf("cryptex: create AES cipher B: %w", aesErr)
		return
	}

	gcmB, aesErr := cipher.NewGCM(blockB)
	if aesErr != nil {
		err = fmt.Errorf("cryptex: create GCM B: %w", aesErr)
		return
	}

	// Encrypt maintext.
	vault_cipher = gcmB.Seal(
		nil,
		nonce_b,
		maintext,
		nil,
	)

	return
}
