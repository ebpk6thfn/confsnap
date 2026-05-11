package encrypt_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/confsnap/internal/encrypt"
)

func TestNew_NotNil(t *testing.T) {
	e := encrypt.New("secret")
	if e == nil {
		t.Fatal("expected non-nil Encryptor")
	}
}

func TestEncrypt_ThenDecrypt_RoundTrip(t *testing.T) {
	e := encrypt.New("my-passphrase")
	plaintext := []byte("hostname = web01\nport = 443")

	ct, err := e.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	got, err := e.Decrypt(ct)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Errorf("round-trip mismatch: got %q, want %q", got, plaintext)
	}
}

func TestEncrypt_ProducesUniqueCiphertexts(t *testing.T) {
	e := encrypt.New("passphrase")
	plaintext := []byte("same input")

	ct1, _ := e.Encrypt(plaintext)
	ct2, _ := e.Encrypt(plaintext)
	if bytes.Equal(ct1, ct2) {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestDecrypt_InvalidData_ReturnsError(t *testing.T) {
	e := encrypt.New("passphrase")
	_, err := e.Decrypt([]byte("not valid ciphertext"))
	if err == nil {
		t.Fatal("expected error for invalid ciphertext")
	}
}

func TestDecrypt_TooShort_ReturnsError(t *testing.T) {
	e := encrypt.New("passphrase")
	_, err := e.Decrypt([]byte{0x01, 0x02})
	if err != encrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_WrongKey_ReturnsError(t *testing.T) {
	e1 := encrypt.New("correct-key")
	e2 := encrypt.New("wrong-key")

	ct, err := e1.Encrypt([]byte("secret data"))
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	_, err = e2.Decrypt(ct)
	if err != encrypt.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestEncrypt_EmptyPlaintext(t *testing.T) {
	e := encrypt.New("passphrase")
	ct, err := e.Encrypt([]byte{})
	if err != nil {
		t.Fatalf("Encrypt empty: %v", err)
	}
	got, err := e.Decrypt(ct)
	if err != nil {
		t.Fatalf("Decrypt empty: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty plaintext, got %q", got)
	}
}
