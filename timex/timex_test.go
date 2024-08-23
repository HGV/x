package timex

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSeries(t *testing.T) {
	start := time.Now()
	stop := start.AddDate(0, 0, 9)

	t.Run("should return nil", func(t *testing.T) {
		assert.Nil(t, GenerateSeries(start, stop, 0))
		assert.Nil(t, GenerateSeries(stop, start, Day))
	})

	t.Run("should return a series of dates", func(t *testing.T) {
		actual := GenerateSeries(start, start.AddDate(0, 0, 9), Day)
		expected := make([]time.Time, 10)
		for i := range expected {
			expected[i] = start.AddDate(0, 0, i)
		}
		assert.True(t, slices.EqualFunc(actual, expected, time.Time.Equal))
	})
}
