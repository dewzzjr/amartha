package action

import "time"

func Filter(start, end *time.Time, dateOnly bool) func(time.Time) bool {
	return func(timestamp time.Time) bool {
		afterStart := true
		matchStart := false
		if start != nil {
			afterStart = timestamp.After(*start)
			matchStart = timestamp.Equal(*start)
		}
		beforeEnd := true
		matchEnd := false
		if end != nil {
			beforeEnd = timestamp.Before(*end)
			matchEnd = timestamp.Equal(*end)
		}
		between := afterStart && beforeEnd
		if !dateOnly {
			return between
		}
		return between || matchStart || matchEnd
	}
}

func DeleteAndPop[T any](slice *[]T, s int) (result T) {
	result = (*slice)[s]
	*slice = append((*slice)[:s], (*slice)[s+1:]...)
	return
}

func CompareTimestampWithDate(timestamp, date time.Time) bool {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location()).Equal(date)
}
