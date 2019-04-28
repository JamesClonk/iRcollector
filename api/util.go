package api

import (
	"strconv"
	"time"
)

type unixTime struct {
	time.Time
}

func (u *unixTime) UnmarshalJSON(data []byte) error {
	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*u = unixTime{time.Unix(unix/1000, 0)}
	return nil
}
