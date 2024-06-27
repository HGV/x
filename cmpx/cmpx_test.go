package cmpx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	tests := []struct {
		a, b     time.Time
		expected func(int) bool
	}{
		{
			a:        time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			b:        time.Date(2024, 6, 1, 11, 0, 0, 0, time.UTC),
			expected: func(i int) bool { return i > 0 },
		},
		{
			a:        time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			b:        time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			expected: func(i int) bool { return i == 0 },
		},
		{
			a:        time.Date(2024, 6, 1, 11, 0, 0, 0, time.UTC),
			b:        time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			expected: func(i int) bool { return i < 0 },
		},
		{
			a:        time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			b:        time.Date(2024, 6, 1, 11, 59, 0, 0, time.UTC),
			expected: func(i int) bool { return i > 0 },
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Time(tt.a, tt.b)
			assert.True(t, tt.expected(result))
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		a, b     bool
		expected func(int) bool
	}{
		{
			a: false, b: true,
			expected: func(i int) bool { return i < 0 },
		},
		{
			a: true, b: false,
			expected: func(i int) bool { return i > 0 },
		},
		{
			a: false, b: false,
			expected: func(i int) bool { return i == 0 },
		},
		{
			a: true, b: true,
			expected: func(i int) bool { return i == 0 },
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Bool(tt.a, tt.b)
			assert.True(t, tt.expected(result))
		})
	}
}
