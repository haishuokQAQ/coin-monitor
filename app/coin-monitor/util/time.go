package util

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Time(ConvertInt642Time(millis))
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	millis := ConvertTime2Int64(time.Time(*t))
	return json.Marshal(millis)
}

func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *Time) Scan(value interface{}) error {
	if v, ok := value.(time.Time); ok {
		*t = Time(v)
		return nil
	} else {
		return errors.Errorf("type convert error : %+v", value)
	}
}

func ConvertTime2Int64(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func ConvertInt642Time(i int64) time.Time {
	msIn1Second := int64(10 * 10 * 10)

	sec := i / msIn1Second
	nsec := (i % msIn1Second) * int64(time.Millisecond)

	result := time.Unix(sec, nsec)

	return result
}
