package timex

import (
	"math/rand"
	"strconv"
	"time"
)

// id长度 = 3位字符（随机数字的36进制） + 8位字符（毫秒时间戳的36进制）
func TimeId() string {
	// 36进制字符范围100-zzz, 对应十进制取值范围1296-46655
	randNum := strconv.FormatInt(int64(rand.Intn(45359)+1296), 36)
	// 36进制8位长度最大值zzzzzzzz, 对应时间戳2821109907455（2059-05-26 01:38:27），该时间之后会变成9位
	timeStr := strconv.FormatInt(time.Now().UnixMilli(), 36)
	return timeStr + randNum
}
