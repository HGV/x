package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeToString(t *testing.T) {
	tests := []struct {
		str  string
		time Time
	}{
		{"13:26:33", Time{13, 26, 33}},
		{"01:02:03", Time{1, 2, 3}},
		{"00:00:00", Time{0, 0, 0}},
	}

	for _, tt := range tests {
		gotTime, err := ParseTime(tt.str)
		assert.Nil(t, err)
		assert.Equal(t, tt.time, gotTime)
		assert.Equal(t, tt.str, tt.time.String())
	}
}

func TestNewTime(t *testing.T) {
	tests := []struct {
		time time.Time
		want Time
	}{
		{time.Date(2014, 8, 20, 15, 8, 43, 0, time.Local), Time{15, 8, 43}},
		{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), Time{0, 0, 0}},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, NewTime(tt.time))
	}
}

func TestTimeIsZero(t *testing.T) {
	tests := []struct {
		time Time
		want bool
	}{
		{Time{0, 0, 0}, true},
		{Time{}, true},
		{Time{0, 0, 1}, false},
		{Time{-1, 0, 0}, false},
		{Time{0, -1, 0}, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.time.IsZero())
	}
}

func TestTimeBefore(t *testing.T) {
	tests := []struct {
		t1, t2 Time
		want   bool
	}{
		{Time{12, 0, 0}, Time{14, 0, 0}, true},
		{Time{12, 20, 0}, Time{12, 30, 0}, true},
		{Time{12, 20, 10}, Time{12, 20, 20}, true},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.t1.Before(tt.t2))
	}
}

func TestTimeAfter(t *testing.T) {
	tests := []struct {
		t1, t2 Time
		want   bool
	}{
		{Time{12, 0, 0}, Time{14, 0, 0}, false},
		{Time{12, 20, 0}, Time{12, 30, 0}, false},
		{Time{12, 20, 10}, Time{12, 20, 20}, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.t1.After(tt.t2))
	}
}

func TestTimeCompare(t *testing.T) {
	tests := []struct {
		t1, t2 Time
		want   int
	}{
		{Time{12, 0, 0}, Time{14, 0, 0}, -1},
		{Time{12, 20, 0}, Time{12, 30, 0}, -1},
		{Time{12, 20, 10}, Time{12, 20, 20}, -1},
		{Time{14, 0, 0}, Time{12, 0, 0}, +1},
		{Time{12, 30, 0}, Time{12, 20, 0}, +1},
		{Time{12, 20, 20}, Time{12, 20, 10}, +1},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.t1.Compare(tt.t2))
	}
}

func TestNewTimeFromMicroseconds(t *testing.T) {
	tests := []struct {
		microseconds int64
		want         Time
	}{
		{0, Time{0, 0, 0}},
		{1_000_000, Time{0, 0, 1}},
		{60_000_000, Time{0, 1, 0}},
		{3600_000_000, Time{1, 0, 0}},
		{3661_000_000, Time{1, 1, 1}},
		{46923_000_000, Time{13, 2, 3}},
		{86399_000_000, Time{23, 59, 59}},
	}

	for _, tt := range tests {
		got := newTimeFromMicroseconds(tt.microseconds)
		assert.Equal(t, tt.want, got)
	}
}

func TestTimeMicroseconds(t *testing.T) {
	tests := []struct {
		time Time
		want int64
	}{
		{Time{0, 0, 0}, 0},
		{Time{0, 0, 1}, 1_000_000},
		{Time{0, 1, 0}, 60_000_000},
		{Time{1, 0, 0}, 3600_000_000},
		{Time{1, 1, 1}, 3661_000_000},
		{Time{13, 2, 3}, 46923_000_000},
		{Time{23, 59, 59}, 86399_000_000},
	}

	for _, tt := range tests {
		got := tt.time.microseconds()
		assert.Equal(t, tt.want, got)
	}
}

func TestMicrosecondsRoundTrip(t *testing.T) {
	tests := []Time{
		{0, 0, 0},
		{1, 2, 3},
		{12, 34, 56},
		{23, 59, 59},
	}

	for _, tt := range tests {
		microseconds := tt.microseconds()
		reconstructed := newTimeFromMicroseconds(microseconds)
		assert.Equal(t, tt, reconstructed)
	}
}
