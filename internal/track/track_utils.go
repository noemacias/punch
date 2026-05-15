package track

import "time"

func LastNDaysRangeStr(days int, layout string) (begin string, end string) {
	if days <= 0 {
		return "", ""
	}

	now := time.Now()
	loc := now.Location()

	beginTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		loc,
	).AddDate(0, 0, -(days - 1))

	endTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		23, 59, 59, int(time.Second-time.Nanosecond),
		loc,
	)

	return beginTime.Format(layout), endTime.Format(layout)
}

func LastNDaysRange(days int, layout string) (begin time.Time, end time.Time) {
	if days <= 0 {
		return time.Time{}, time.Time{}
	}

	now := time.Now()
	loc := now.Location()

	beginTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		loc,
	).AddDate(0, 0, -(days - 1))

	endTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		23, 59, 59, int(time.Second-time.Nanosecond),
		loc,
	)

	return beginTime, endTime
}
