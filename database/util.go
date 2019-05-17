package database

import (
	"fmt"
	"time"
)

type Laptime int

func (l Laptime) String() string {
	return fmt.Sprintf("%s", time.Duration(l*100)*time.Microsecond)
}
