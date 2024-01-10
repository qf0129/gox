package securex

import (
	"encoding/binary"
	"encoding/hex"
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

	token, err := DesEncrypt(tokenByte, []byte(secretKey))
	return hex.EncodeToString(token), err
}

// 解析令牌
func ParseToken(token string, secretKey string, expiredSeconds int64) (string, error) {
	text, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}
	tokenStr, err := DesDecrypt([]byte(text), []byte(secretKey))
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
