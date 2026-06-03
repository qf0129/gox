package timex

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/qf0129/gox/pkg/ptrx"
)

type Time time.Time

const DefaultFormat = "2006-01-02 15:04:05"

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	str := string(data)
	t1, err := time.Parse(DefaultFormat, strings.Trim(str, "\""))
	*t = Time(t1)
	return err
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", time.Time(*t).Format(DefaultFormat))), nil
}

var zeroTime time.Time

func (t Time) Value() (driver.Value, error) {
	tTime := time.Time(t)
	if tTime.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tTime.Format(DefaultFormat), nil
}

func (t *Time) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timex.Time", v)
}

func (t *Time) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

func (t *Time) Time() time.Time {
	return time.Time(*t)
}

func Now() Time {
	return Time(time.Now())
}

func NowP() *Time {
	return ptrx.Ptr(Now())
}
