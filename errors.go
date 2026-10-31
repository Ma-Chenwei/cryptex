package cryptex

import "errors"

var (
	ErrInvalidPassword = errors.New("cryptex: invalid password")

	ErrInvalidKey = errors.New("cryptex: invalid system key")

	ErrInvalidSalt = errors.New("cryptex: invalid salt")

	ErrInvalidNonce = errors.New("cryptex: invalid nonce")

	ErrInvalidCiphertext = errors.New("cryptex: invalid ciphertext")

	ErrRandomGeneration = errors.New("cryptex: random generation failed")
)
