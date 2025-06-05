package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtrEqual(t *testing.T) {
	i1, i2, i3 := 10, 10, 20

	tests := []struct {
		name     string
		a        *int
		b        *int
		expected bool
	}{
		{"both nil", nil, nil, true},
		{"a nil, b non-nil", nil, &i1, false},
		{"a non-nil, b nil", &i1, nil, false},
		{"same value", &i1, &i2, true},
		{"different value", &i1, &i3, false},
		{"same pointer", &i1, &i1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, PtrEqual(tt.a, tt.b))
		})
	}
}
