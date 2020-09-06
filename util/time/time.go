package time

import (
	"time"
)

// Constants of time
const (
	Year   string = "2006"
	Month  string = "01"
	Day    string = "02"
	Hour   string = "15"
	Minute string = "04"
	Second string = "05"

	DateSlashFormat     = Year + "/" + Month + "/" + Day
	FullTimeFormat = Year + "-" + Month + "-" + Day + " " + Hour + ":" + Minute + ":" + Second
	QueryFormat            = Year + "-" + Month + "-" + Day
)

// GetDateSlashFormatTime to get the time format for DateSlashFormat
func GetDateSlashFormatTime(t time.Time) string {
	return t.Local().Format(DateSlashFormat)
}

// GetFullTimeFormat to get the full time format
func GetFullTimeFormat(t time.Time) string {
	return t.Local().Format(FullTimeFormat)
}

// ConvertFullTimeFormat to convert FullTimeFormat str to time pointer
func ConvertFullTimeFormat(t string) *time.Time {
	return convertTime(t, FullTimeFormat)
}

// ConvertStrToTime to convert string to time pointer
func ConvertStrToTime(t string) *time.Time {
	return convertTime(t, QueryFormat)
}

func convertTime(t string, format string) *time.Time {
	var tp *time.Time
	if t != "" {
		nt, _ := time.ParseInLocation(format, t, time.Local)
		tp = &nt
	} else {
		tp = nil
	}
	return tp
}
