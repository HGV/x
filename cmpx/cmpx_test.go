package cmpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
