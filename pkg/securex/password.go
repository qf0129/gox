package securex

import "golang.org/x/crypto/bcrypt"

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
