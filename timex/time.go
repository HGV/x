package timex

import (
	"encoding"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Time struct {
	Hour   int
	Minute int
	Second int
}

func NewTime(t time.Time) Time {
	var tm Time
	tm.Hour, tm.Minute, tm.Second = t.Clock()
	return tm
}

func ParseTime(s string) (Time, error) {
	t, err := time.Parse(time.TimeOnly, s)
	if err != nil {
		return Time{}, err
	}
	return NewTime(t), nil
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

func (t Time) IsZero() bool {
	return t.Hour == 0 && t.Minute == 0 && t.Second == 0
}

func (t Time) Before(t2 Time) bool {
	if t.Hour != t2.Hour {
		return t.Hour < t2.Hour
	}
	if t.Minute != t2.Minute {
		return t.Minute < t2.Minute
	}
	return t.Second < t2.Second
}

func (t Time) After(t2 Time) bool {
	return t2.Before(t)
}

func (t Time) Compare(t2 Time) int {
	if t.Before(t2) {
		return -1
	} else if t.After(t2) {
		return +1
	}
	return 0
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Time) UnmarshalText(data []byte) error {
	var err error
	*t, err = ParseTime(string(data))
	return err
}

var _ encoding.TextMarshaler = Time{}
var _ encoding.TextUnmarshaler = &Time{}

func (t *Time) ScanTime(v pgtype.Time) error {
	if !v.Valid {
		*t = Time{}
		return nil
	}

	*t = newTimeFromMicroseconds(v.Microseconds)
	return nil
}

func (t Time) TimeValue() (pgtype.Time, error) {
	return pgtype.Time{
		Microseconds: t.microseconds(),
		Valid:        true,
	}, nil
}

var _ pgtype.TimeScanner = &Time{}
var _ pgtype.TimeValuer = Time{}

func newTimeFromMicroseconds(usec int64) Time {
	totalSeconds := usec / 1_000_000
	return Time{
		Hour:   int(totalSeconds / 3600),
		Minute: int((totalSeconds % 3600) / 60),
		Second: int(totalSeconds % 60),
	}
}

func (t Time) microseconds() int64 {
	return int64((t.Hour*3600 + t.Minute*60 + t.Second) * 1_000_000)
}
