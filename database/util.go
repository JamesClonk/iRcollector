package database

import (
	"fmt"
	"time"
)

type Laptime int

func (l Laptime) String() string {
	if l == 0 {
		return ""
	}
	return fmt.Sprintf("%s", time.Duration(l*100)*time.Microsecond)
}

func (l Laptime) Seconds() int {
	if l == 0 {
		return 0
	}
	return int((time.Duration(l*100) * time.Microsecond).Seconds())
}

func (l Laptime) Milliseconds() int64 {
	if l == 0 {
		return 0
	}
	return (time.Duration(l*100) * time.Microsecond).Milliseconds()
}

func WeekStart(reference time.Time) time.Time {
	y, m, d := reference.Date()

	t := time.Date(y, m, d, 0, 0, 0, 0, reference.Location())
	weekday := int(t.Weekday())

	weekStartDayInt := int(time.Tuesday)
	if weekday < weekStartDayInt {
		weekday = weekday + 7 - weekStartDayInt
	} else {
		weekday = weekday - weekStartDayInt
	}
	return t.AddDate(0, 0, -weekday)
}
