package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

type encodedTime struct {
	Time time.Time
}

func (e *encodedTime) UnmarshalJSON(data []byte) error {
	input := strings.Replace(decode(string(data)), `"`, "", -1)
	if strings.Count(input, ":") == 1 {
		input = input + ":00"
	}

	t, err := time.Parse("2006-01-02 15:04:05", input)
	if err != nil {
		return err
	}

	*e = encodedTime{t}
	return nil
}

type encodedString string

func (e encodedString) String() string {
	return decode(string(e))
}

func (e encodedString) Laptime() int {
	input := e.String()
	if len(input) == 0 {
		return -1
	}
	input = strings.Replace(input, ":", "m", -1)
	input = strings.Replace(input, ".", "s", -1)
	input = input + "ms"

	d, err := time.ParseDuration(input)
	if err != nil {
		return -1
	}
	return int(d.Nanoseconds() / 1000 / 100)
}

func decode(value string) string {
	decodedValue, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return decodedValue
}

type laptime int

func (l laptime) String() string {
	if l == 0 {
		return ""
	}
	return fmt.Sprintf("%s", time.Duration(l*100)*time.Microsecond)
}
