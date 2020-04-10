package kreedz

import (
	"fmt"
	"time"
)

func SecondsToMinutes(secs float64) string {
	min := int64(secs / 60)
	sec := secs - float64(min) * 60
	zero := ""
	if sec < 10 {
		zero = "0"
	}
	return fmt.Sprintf("%02d:%s%.2f", min, zero, sec)
}

func FormatDate(t time.Time) string {
	return t.UTC().Format("2006-01-02")
}