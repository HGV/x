package timex

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type DateRange struct {
	Start Date
	End   Date
}

func (d *DateRange) ScanNull() error {
	return fmt.Errorf("cannot scan NULL into daterange")
}

func (d *DateRange) ScanBounds() (lowerTarget, upperTarget any) {
	return &d.Start, &d.End
}

func (d *DateRange) SetBoundTypes(lower, upper pgtype.BoundType) error {
	if lower == pgtype.Exclusive {
		d.Start = d.Start.AddDays(1)
	}
	if upper == pgtype.Exclusive {
		d.End = d.End.AddDays(-1)
	}
	return nil
}

func (d DateRange) IsNull() bool {
	return false
}

func (d DateRange) BoundTypes() (lower, upper pgtype.BoundType) {
	return pgtype.Inclusive, pgtype.Inclusive
}

func (d DateRange) Bounds() (lower, upper any) {
	return d.Start, d.End
}

var _ pgtype.RangeScanner = &DateRange{}
var _ pgtype.RangeValuer = DateRange{}

type NullDateRange struct {
	DateRange DateRange
	Valid     bool
}

func (nd *NullDateRange) ScanNull() error {
	*nd = NullDateRange{}
	return nil
}

func (nd *NullDateRange) ScanBounds() (lowerTarget, upperTarget any) {
	return nd.DateRange.ScanBounds()
}

func (nd *NullDateRange) SetBoundTypes(lower, upper pgtype.BoundType) error {
	nd.DateRange.SetBoundTypes(lower, upper)
	nd.Valid = true
	return nil
}

func (nd NullDateRange) IsNull() bool {
	return !nd.Valid
}

func (nd NullDateRange) BoundTypes() (lower, upper pgtype.BoundType) {
	return nd.DateRange.BoundTypes()
}

func (nd NullDateRange) Bounds() (lower, upper any) {
	return nd.DateRange.Bounds()
}

var _ pgtype.RangeScanner = &NullDateRange{}
var _ pgtype.RangeValuer = NullDateRange{}
