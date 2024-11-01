package timex

import (
	"encoding"
	"fmt"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func NewDateFromTime(t time.Time) Date {
	var d Date
	d.Year, d.Month, d.Day = t.Date()
	return d
}

func ParseDate(s string) (Date, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return Date{}, err
	}
	return NewDateFromTime(t), nil
}

func (d Date) Before(d2 Date) bool {
	if d.Year != d2.Year {
		return d.Year < d2.Year
	}
	if d.Month != d2.Month {
		return d.Month < d2.Month
	}
	return d.Day < d2.Day
}

func (d Date) After(d2 Date) bool {
	return d2.Before(d)
}

func (d Date) Compare(d2 Date) int {
	if d.Before(d2) {
		return -1
	}
	if d.After(d2) {
		return +1
	}
	return 0
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Date) UnmarshalText(data []byte) error {
	var err error
	*d, err = ParseDate(string(data))
	return err
}

var _ encoding.TextMarshaler = Date{}
var _ encoding.TextUnmarshaler = &Date{}
