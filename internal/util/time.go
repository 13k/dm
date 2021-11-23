package util

import (
	"time"
)

func Today() string {
	return time.Now().Format("2006-01-02")
}
