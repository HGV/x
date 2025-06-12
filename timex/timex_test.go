package timex

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeTimeSeries(t *testing.T) {
	start := time.Now()
	stop := start.AddDate(0, 0, 9)

	t.Run("should return nil", func(t *testing.T) {
		assert.Nil(t, MakeTimeSeries(start, stop, 0))
		assert.Nil(t, MakeTimeSeries(stop, start, Day))
	})

	t.Run("should return a series of dates", func(t *testing.T) {
		actual := MakeTimeSeries(start, start.AddDate(0, 0, 9), Day)
		expected := make([]time.Time, 10)
		for i := range expected {
			expected[i] = start.AddDate(0, 0, i)
		}
		assert.True(t, slices.EqualFunc(actual, expected, time.Time.Equal))
	})
}

func TestMakeDateSeries(t *testing.T) {
	start := Today()
	stop := start.AddDays(9)

	t.Run("should return nil", func(t *testing.T) {
		assert.Nil(t, MakeDateSeries(start, stop, 0))
		assert.Nil(t, MakeDateSeries(stop, start, 1))
	})

	t.Run("should return a series of dates", func(t *testing.T) {
		actual := MakeDateSeries(start, start.AddDays(9), 1)
		expected := make([]Date, 10)
		for i := range expected {
			expected[i] = start.AddDays(i)
		}
		assert.True(t, slices.EqualFunc(actual, expected, func(a, b Date) bool {
			return a.Compare(b) == 0
		}))
	})
}
