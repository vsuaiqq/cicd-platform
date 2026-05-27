package service

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func testKeyBase64(t *testing.T) string {
	t.Helper()
	raw := bytes.Repeat([]byte("k"), 32)
	return base64.StdEncoding.EncodeToString(raw)
}

func TestEncrypter_roundTrip(t *testing.T) {
	enc, err := NewEncrypter(testKeyBase64(t))
	if err != nil {
		t.Fatalf("NewEncrypter: %v", err)
	}

	plain := []byte("ssh-private-key-material")
	cipher, err := enc.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if bytes.Equal(cipher, plain) {
		t.Fatal("ciphertext equals plaintext")
	}

	got, err := enc.Decrypt(cipher)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("decrypted = %q", got)
	}
}

func TestEncrypter_rejectsShortKey(t *testing.T) {
	key := base64.StdEncoding.EncodeToString([]byte("short"))
	if _, err := NewEncrypter(key); err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestEncrypter_rejectsTamperedCiphertext(t *testing.T) {
	enc, err := NewEncrypter(testKeyBase64(t))
	if err != nil {
		t.Fatalf("NewEncrypter: %v", err)
	}
	cipher, err := enc.Encrypt([]byte("secret"))
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	cipher[len(cipher)-1] ^= 0xff
	if _, err := enc.Decrypt(cipher); err == nil {
		t.Fatal("expected decrypt error")
	}
}
