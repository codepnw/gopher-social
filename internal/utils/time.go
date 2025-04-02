package utils

import "time"

func TimeString() string {
	format := "2006-01-02 15:04:05"
	return time.Now().Format(format)
}