package securex

import (
	"encoding/binary"
	"errors"
	"time"
)

var ErrTokenExpired = errors.New("TokenExpired")

// 创建令牌
func CreateToken(body string, secretKey string) (string, error) {
	tokenByte := make([]byte, len(body)+8) // 前8位放时间戳
	binary.BigEndian.PutUint64(tokenByte[:8], uint64(time.Now().Unix()))
	for i, v := range body {
		tokenByte[i+8] = byte(v)
	}

	return Encrypt(tokenByte, []byte(secretKey))
}

// 解析令牌
func ParseToken(token string, secretKey string, expiredSeconds int64) (string, error) {
	tokenStr, err := Decrypt(token, []byte(secretKey))
	if err != nil {
		return "", err
	}

	tokenTs := int64(binary.BigEndian.Uint64(tokenStr[0:8]))
	nowTs := time.Now().Unix()

	if nowTs > tokenTs+expiredSeconds {
		return "", ErrTokenExpired
	}

	body := string(tokenStr[8:])
	return body, nil
}
