package envfile

import (
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "super_secret_value"
	passphrase := "my-passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if encrypted == plaintext {
		t.Fatal("expected encrypted to differ from plaintext")
	}

	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	encrypted, err := Encrypt("value", "correct")
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	_, err = Decrypt(encrypted, "wrong")
	if err == nil {
		t.Fatal("expected error decrypting with wrong passphrase")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("!!!notbase64!!!", "pass")
	if err == nil {
		t.Fatal("expected error on invalid base64")
	}
}

func TestEncrypt_Nondeterministic(t *testing.T) {
	a, _ := Encrypt("value", "pass")
	b, _ := Encrypt("value", "pass")
	if a == b {
		t.Error("expected different ciphertexts due to random nonce")
	}
}
