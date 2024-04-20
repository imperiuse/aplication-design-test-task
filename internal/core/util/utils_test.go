package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSliceForDuplicates(t *testing.T) {
	tests := []struct {
		name        string
		input       []int
		expectError bool
	}{
		{"No duplicates", []int{1, 2, 3, 4, 5}, false},
		{"With duplicates", []int{1, 2, 2, 4, 5}, true},
		{"Empty slice", []int{}, false},
		{"Single element", []int{1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSliceForDuplicates(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsDayBetween(t *testing.T) {
	from := NewDay(2024, 1, 1)
	to := NewDay(2024, 1, 10)
	dayInside := NewDay(2024, 1, 5)
	dayOutside := NewDay(2024, 1, 15)

	assert.True(t, IsDayBetween(dayInside, from, to))
	assert.False(t, IsDayBetween(dayOutside, from, to))
}

func TestDaysBetween(t *testing.T) {
	from := NewDay(2024, 1, 1)
	to := NewDay(2024, 1, 3)
	expectedDays := []Day{
		NewDay(2024, 1, 1),
		NewDay(2024, 1, 2),
		NewDay(2024, 1, 3),
	}

	result := DaysBetween(from, to)
	assert.Equal(t, expectedDays, result)
}

func TestDaysBetween_EmptySliceWhenFromAfterTo(t *testing.T) {
	from := NewDay(2024, 1, 3)
	to := NewDay(2024, 1, 1)
	assert.Empty(t, DaysBetween(from, to))
}
