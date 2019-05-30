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
