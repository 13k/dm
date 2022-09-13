package util

import (
	"time"
)

func TodayString() string {
	return time.Now().Format("2006-01-02")
}
