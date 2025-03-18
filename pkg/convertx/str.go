package convertx

import "strconv"

func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
