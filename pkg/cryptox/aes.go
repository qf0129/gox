package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	AES256GCMAlgorithm     = "AES-256-GCM"
	EncryptedValueVersion  = 1
	credentialVersionToken = "v1"
	credentialPrefix       = credentialVersionToken + ":"
)

var (
	ErrEmptyCredentialSecret     = errors.New("credential secret is empty")
	ErrInvalidAES256GCMKeyLength = errors.New("aes-256-gcm key must be 32 bytes")
	ErrInvalidCredentialFormat   = errors.New("invalid credential ciphertext format")
	ErrUnsupportedEncryptedValue = errors.New("unsupported encrypted value")
)

type EncryptedValue struct {
	Alg        string `json:"alg"`
	Version    int    `json:"version"`
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

func EncryptAES256GCM(plaintext string, key []byte) (*EncryptedValue, error) {
	if len(key) != 32 {
		return nil, ErrInvalidAES256GCMKeyLength
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	return &EncryptedValue{
		Alg:        AES256GCMAlgorithm,
		Version:    EncryptedValueVersion,
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}, nil
}

func DecryptAES256GCM(value EncryptedValue, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidAES256GCMKeyLength
	}
	if value.Alg != "" && value.Alg != AES256GCMAlgorithm {
		return "", fmt.Errorf("%w: alg=%s", ErrUnsupportedEncryptedValue, value.Alg)
	}
	if value.Version != 0 && value.Version != EncryptedValueVersion {
		return "", fmt.Errorf("%w: version=%d", ErrUnsupportedEncryptedValue, value.Version)
	}

	nonce, err := base64.StdEncoding.DecodeString(value.Nonce)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(value.Ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func EncryptCredential(plaintext string, secret string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key, err := DeriveCredentialKey(secret)
	if err != nil {
		return "", err
	}

	value, err := EncryptAES256GCM(plaintext, key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s:%s", credentialPrefix, value.Nonce, value.Ciphertext), nil
}

func DecryptCredential(stored string, secret string) (string, error) {
	if stored == "" {
		return "", nil
	}
	if !strings.HasPrefix(stored, credentialPrefix) {
		return stored, nil
	}

	encodedValue, ok := strings.CutPrefix(stored, credentialPrefix)
	if !ok {
		return "", ErrInvalidCredentialFormat
	}
	nonce, ciphertext, ok := strings.Cut(encodedValue, ":")
	if !ok || nonce == "" || ciphertext == "" {
		return "", ErrInvalidCredentialFormat
	}

	key, err := DeriveCredentialKey(secret)
	if err != nil {
		return "", err
	}

	return DecryptAES256GCM(EncryptedValue{
		Alg:        AES256GCMAlgorithm,
		Version:    EncryptedValueVersion,
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}, key)
}

func DeriveCredentialKey(secret string) ([]byte, error) {
	if secret == "" {
		return nil, ErrEmptyCredentialSecret
	}

	if decoded, err := base64.StdEncoding.DecodeString(secret); err == nil && len(decoded) == 32 {
		return decoded, nil
	}

	if len([]byte(secret)) == 32 {
		return []byte(secret), nil
	}

	sum := sha256.Sum256([]byte(secret))
	return sum[:], nil
}
