package hashx

import (
	"crypto/md5"
	"encoding/hex"
)

func GetStrMD5(s string) string {
	return GetMD5Hash([]byte(s))
}

func GetMD5Hash(b []byte) string {
	binHash := md5.Sum(b)
	return hex.EncodeToString(binHash[:])
}
