package util

import (
	"fmt"
	"time"
)

// CheckSliceForDuplicates ensures that all value in the slice are unique.
func CheckSliceForDuplicates[T comparable](s []T) error {
	seen := make(map[T]struct{})
	for _, v := range s {
		if _, exists := seen[v]; exists {
			return fmt.Errorf("duplicate value found: %v", v)
		}
		seen[v] = struct{}{}
	}

	return nil
}

type (
	Day  = time.Time
	Days = []Day
)

func IsDayBetween(day Day, from Day, to Day) bool {
	day = ToDay(day)
	return !day.Before(ToDay(from)) && !day.After(ToDay(to))
}

func DaysBetween(from time.Time, to time.Time) Days {
	if from.After(to) {
		return nil
	}

	days := make([]time.Time, 0)
	for d := ToDay(from); !d.After(ToDay(to)); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}

func ToDay(timestamp time.Time) Day {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}

func NewDay(year int, month int, day int) Day {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
