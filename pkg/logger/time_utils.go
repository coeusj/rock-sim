package logger

import "time"

func DateTimeWithNanoseconds(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.9999999")
}
