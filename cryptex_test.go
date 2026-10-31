package cryptex

import (
	"bytes"
	"testing"
)

func TestGenerateC(t *testing.T) {
	c, err := GenerateC()

	if err != nil {
		t.Fatal(err)
	}

	if len(c) != CSize {
		t.Fatalf("expected %d bytes, got %d", CSize, len(c))
	}
}

func TestEncryptDecryptAutoC(t *testing.T) {
	password := "MyMasterPassword"
	original := []byte("Hello Cryptex")

	c, sa, na, sc, sb, nb, vc, err :=
		Encrypt(password, nil, original)

	if err != nil {
		t.Fatal(err)
	}

	if len(c) != CSize {
		t.Fatal("invalid generated C")
	}

	result, err := Decrypt(
		password,
		sa, na, sc,
		sb, nb, vc,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(original, result) {
		t.Fatal("data mismatch")
	}
}

func TestEncryptDecryptEmptyC(t *testing.T) {
	password := "MyMasterPassword"
	original := []byte("Empty C test")

	c, sa, na, sc, sb, nb, vc, err :=
		Encrypt(password, []byte{}, original)

	if err != nil {
		t.Fatal(err)
	}

	if len(c) != CSize {
		t.Fatal("empty C should generate a new C")
	}

	result, err := Decrypt(
		password,
		sa, na, sc,
		sb, nb, vc,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(original, result) {
		t.Fatal("data mismatch")
	}
}

func TestEncryptDecryptCustomC(t *testing.T) {
	password := "MyMasterPassword"
	original := []byte("Custom C test")

	c := bytes.Repeat([]byte{0x42}, CSize)

	actualC, sa, na, sc, sb, nb, vc, err :=
		Encrypt(password, c, original)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(c, actualC) {
		t.Fatal("returned C does not match input C")
	}

	result, err := Decrypt(
		password,
		sa, na, sc,
		sb, nb, vc,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(original, result) {
		t.Fatal("data mismatch")
	}
}

func TestInvalidC(t *testing.T) {
	password := "MyMasterPassword"

	invalid := []byte{1, 2, 3}

	_, _, _, _, _, _, _, err :=
		Encrypt(password, invalid, []byte("test"))

	if err != ErrInvalidKey {
		t.Fatalf("expected ErrInvalidKey, got %v", err)
	}
}

func TestEmptyPassword(t *testing.T) {
	_, _, _, _, _, _, _, err :=
		Encrypt("", nil, []byte("test"))

	if err != ErrInvalidPassword {
		t.Fatal("expected ErrInvalidPassword")
	}
}

func TestWrongPassword(t *testing.T) {
	c, sa, na, sc, sb, nb, vc, err :=
		Encrypt("password1", nil, []byte("secret"))

	_ = c

	if err != nil {
		t.Fatal(err)
	}

	_, err = Decrypt(
		"password2",
		sa, na, sc,
		sb, nb, vc,
	)

	if err == nil {
		t.Fatal("wrong password should fail")
	}
}

func TestTamperedCiphertext(t *testing.T) {
	password := "MyMasterPassword"

	_, sa, na, sc, sb, nb, vc, err :=
		Encrypt(password, nil, []byte("secret"))

	if err != nil {
		t.Fatal(err)
	}

	vc[0] ^= 0x01

	_, err = Decrypt(
		password,
		sa, na, sc,
		sb, nb, vc,
	)

	if err == nil {
		t.Fatal("tampered ciphertext should fail")
	}
}

func TestTamperedSystemCipher(t *testing.T) {
	password := "MyMasterPassword"

	_, sa, na, sc, sb, nb, vc, err :=
		Encrypt(password, nil, []byte("secret"))

	if err != nil {
		t.Fatal(err)
	}

	sc[0] ^= 0x01

	_, err = Decrypt(
		password,
		sa, na, sc,
		sb, nb, vc,
	)

	if err == nil {
		t.Fatal("tampered system cipher should fail")
	}
}

func TestEmptyMaintext(t *testing.T) {
	password := "MyMasterPassword"

	_, sa, na, sc, sb, nb, vc, err :=
		Encrypt(password, nil, []byte{})

	if err != nil {
		t.Fatal(err)
	}

	result, err := Decrypt(
		password,
		sa, na, sc,
		sb, nb, vc,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatal("expected empty plaintext")
	}
}

func TestLargeData(t *testing.T) {
	password := "MyMasterPassword"

	original := bytes.Repeat([]byte("Cryptex"), 1024*1024)

	_, sa, na, sc, sb, nb, vc, err :=
		Encrypt(password, nil, original)

	if err != nil {
		t.Fatal(err)
	}

	result, err := Decrypt(
		password,
		sa, na, sc,
		sb, nb, vc,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(original, result) {
		t.Fatal("large data mismatch")
	}
}

func TestRandomness(t *testing.T) {
	password := "MyMasterPassword"
	original := []byte("random test")

	_, sa1, na1, _, sb1, nb1, _, err :=
		Encrypt(password, nil, original)

	if err != nil {
		t.Fatal(err)
	}

	_, sa2, na2, _, sb2, nb2, _, err :=
		Encrypt(password, nil, original)

	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(sa1, sa2) {
		t.Fatal("salt_a repeated")
	}

	if bytes.Equal(sb1, sb2) {
		t.Fatal("salt_b repeated")
	}

	if bytes.Equal(na1, na2) {
		t.Fatal("nonce_a repeated")
	}

	if bytes.Equal(nb1, nb2) {
		t.Fatal("nonce_b repeated")
	}
}

func BenchmarkEncrypt(b *testing.B) {
	password := "MyMasterPassword"
	data := []byte("Hello Cryptex")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _, _, _, _, _, err :=
			Encrypt(password, nil, data)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	password := "MyMasterPassword"
	data := []byte("Hello Cryptex")

	_, sa, na, sc, sb, nb, vc, err :=
		Encrypt(password, nil, data)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Decrypt(
			password,
			sa, na, sc,
			sb, nb, vc,
		)

		if err != nil {
			b.Fatal(err)
		}
	}
}
