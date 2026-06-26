package cryptox

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestEncryptAES256GCMDecryptAES256GCM(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := "database-password"

	value, err := EncryptAES256GCM(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAES256GCM() error = %v", err)
	}
	if value == nil {
		t.Fatal("EncryptAES256GCM() returned nil value")
	}
	if value.Alg != AES256GCMAlgorithm {
		t.Fatalf("Alg = %q, want %q", value.Alg, AES256GCMAlgorithm)
	}
	if value.Version != EncryptedValueVersion {
		t.Fatalf("Version = %d, want %d", value.Version, EncryptedValueVersion)
	}
	if value.Ciphertext == base64.StdEncoding.EncodeToString([]byte(plaintext)) {
		t.Fatal("ciphertext should not be plaintext encoded as base64")
	}

	got, err := DecryptAES256GCM(*value, key)
	if err != nil {
		t.Fatalf("DecryptAES256GCM() error = %v", err)
	}
	if got != plaintext {
		t.Fatalf("DecryptAES256GCM() = %q, want %q", got, plaintext)
	}
}

func TestEncryptAES256GCMRejectsInvalidKeyLength(t *testing.T) {
	_, err := EncryptAES256GCM("secret", []byte("short-key"))
	if !errors.Is(err, ErrInvalidAES256GCMKeyLength) {
		t.Fatalf("EncryptAES256GCM() error = %v, want %v", err, ErrInvalidAES256GCMKeyLength)
	}
}

func TestDecryptAES256GCMRejectsUnsupportedMetadata(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	value, err := EncryptAES256GCM("secret", key)
	if err != nil {
		t.Fatalf("EncryptAES256GCM() error = %v", err)
	}

	value.Alg = "AES-128-GCM"
	_, err = DecryptAES256GCM(*value, key)
	if !errors.Is(err, ErrUnsupportedEncryptedValue) {
		t.Fatalf("DecryptAES256GCM() error = %v, want %v", err, ErrUnsupportedEncryptedValue)
	}
}

func TestEncryptCredentialDecryptCredential(t *testing.T) {
	const (
		secret    = "x"
		plaintext = "token-value"
	)

	stored, err := EncryptCredential(plaintext, secret)
	if err != nil {
		t.Fatalf("EncryptCredential() error = %v", err)
	}
	if !strings.HasPrefix(stored, credentialPrefix) {
		t.Fatalf("EncryptCredential() = %q, want prefix %q", stored, credentialPrefix)
	}
	if strings.Contains(stored, plaintext) {
		t.Fatalf("EncryptCredential() = %q, should not contain plaintext", stored)
	}

	got, err := DecryptCredential(stored, secret)
	if err != nil {
		t.Fatalf("DecryptCredential() error = %v", err)
	}
	if got != plaintext {
		t.Fatalf("DecryptCredential() = %q, want %q", got, plaintext)
	}
}

func TestDecryptCredentialKeepsLegacyPlaintext(t *testing.T) {
	got, err := DecryptCredential("legacy-token", "secret")
	if err != nil {
		t.Fatalf("DecryptCredential() error = %v", err)
	}
	if got != "legacy-token" {
		t.Fatalf("DecryptCredential() = %q, want legacy plaintext", got)
	}
}

func TestDecryptCredentialRejectsInvalidFormat(t *testing.T) {
	_, err := DecryptCredential("v1:nonce-only", "secret")
	if !errors.Is(err, ErrInvalidCredentialFormat) {
		t.Fatalf("DecryptCredential() error = %v, want %v", err, ErrInvalidCredentialFormat)
	}
}

func TestDeriveCredentialKey(t *testing.T) {
	rawKey := []byte("12345678901234567890123456789012")
	encodedKey := base64.StdEncoding.EncodeToString(rawKey)

	tests := []struct {
		name    string
		secret  string
		want    []byte
		wantLen int
	}{
		{name: "base64 encoded 32 byte key", secret: encodedKey, want: rawKey},
		{name: "raw 32 byte key", secret: string(rawKey), want: rawKey},
		{name: "passphrase", secret: "passphrase", wantLen: 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeriveCredentialKey(tt.secret)
			if err != nil {
				t.Fatalf("DeriveCredentialKey() error = %v", err)
			}
			if tt.want != nil && !bytes.Equal(got, tt.want) {
				t.Fatalf("DeriveCredentialKey() = %q, want %q", got, tt.want)
			}
			if tt.wantLen != 0 && len(got) != tt.wantLen {
				t.Fatalf("DeriveCredentialKey() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestDeriveCredentialKeyRejectsEmptySecret(t *testing.T) {
	_, err := DeriveCredentialKey("")
	if !errors.Is(err, ErrEmptyCredentialSecret) {
		t.Fatalf("DeriveCredentialKey() error = %v, want %v", err, ErrEmptyCredentialSecret)
	}
}

func ExampleEncryptCredential() {
	stored, err := EncryptCredential("api-token", "app-secret")
	if err != nil {
		panic(err)
	}

	plaintext, err := DecryptCredential(stored, "app-secret")
	if err != nil {
		panic(err)
	}

	fmt.Println(plaintext)
	// Output:
	// api-token
}
