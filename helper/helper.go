package helper

import "time"

func FormatTimeRFC3339ToDatetime(timeStr string) string {
	tt, _ := time.ParseInLocation(time.RFC3339, timeStr, time.Local)
	return tt.Format(time.DateTime)
}
