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
