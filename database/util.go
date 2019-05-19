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
