package securex

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// 哈希密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// 校验密码
func VerifyPassword(psd string, psdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(psdHash), []byte(psd))
	return err == nil
}

// 使用CBC模式加密
func Encrypt(text, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// CBC模式需要一个填充器
	padding := aes.BlockSize - len(text)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	text = append(text, padText...)

	// 创建CBC模式的加密器
	mode := cipher.NewCBCEncrypter(block, key)
	ciphertext := make([]byte, len(text))
	mode.CryptBlocks(ciphertext, []byte(text))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// 使用CBC模式解密
func Decrypt(encryptedText string, key []byte) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建CBC模式的解密器，不设置IV向量（使用密钥作为IV）
	mode := cipher.NewCBCDecrypter(block, key)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	padding := int(plaintext[len(plaintext)-1])
	if padding < 1 || int(padding) > len(plaintext) {
		return nil, fmt.Errorf("无效的填充")
	}
	plaintext = plaintext[:len(plaintext)-padding]

	return plaintext, nil
}
