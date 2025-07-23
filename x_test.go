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

func TestIsHGVEmployee(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"user@hgv.it", true},
		{"user@hgj.it", true},
		{"User@HGV.it", true}, // case-insensitive check
		{"User@HGJ.it", true}, // case-insensitive check
		{"user@example.com", false},
		{"@hgv.it", true},
		{"someone@fakehgv.it", false}, // suffix match only
		{"employee@hgv.com", false},
		{"employee@hgv.its", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsHGVEmployee(tt.email))
		})
	}
}
