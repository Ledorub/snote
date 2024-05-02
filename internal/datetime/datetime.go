package datetime

import "time"

func TimeAsLocalTime(t time.Time, tz *time.Location) time.Time {
	t = t.UTC()
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		tz,
	)
}
