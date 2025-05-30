package timex

import (
	"cmp"
	"errors"

	"github.com/HGV/x/cmpx"
	"github.com/jackc/pgx/v5/pgtype"
)

type DaysOfWeek struct {
	Mo, Tu, We, Th, Fr, Sa, Su bool
}

func ParseDaysOfWeek(s []bool) (*DaysOfWeek, error) {
	if len(s) != 7 {
		return nil, errors.New("invalid days of week length")
	}

	return &DaysOfWeek{
		Mo: s[0],
		Tu: s[1],
		We: s[2],
		Th: s[3],
		Fr: s[4],
		Sa: s[5],
		Su: s[6],
	}, nil
}

func (w DaysOfWeek) DaysOfWeek() []bool {
	return []bool{w.Mo, w.Tu, w.We, w.Th, w.Fr, w.Sa, w.Su}
}

func (w DaysOfWeek) Compare(w2 DaysOfWeek) int {
	return cmp.Or(
		cmpx.Bool(w.Mo, w2.Mo),
		cmpx.Bool(w.Tu, w2.Tu),
		cmpx.Bool(w.We, w2.We),
		cmpx.Bool(w.Th, w2.Th),
		cmpx.Bool(w.Fr, w2.Fr),
		cmpx.Bool(w.Sa, w2.Sa),
		cmpx.Bool(w.Su, w2.Su),
	)
}

func (w *DaysOfWeek) ScanBits(bits pgtype.Bits) error {
	bitsLen := int(bits.Len)
	b := bits.Bytes[0]
	s := make([]bool, 0, bitsLen)
	for i := range bitsLen {
		bit := b & (1 << (bitsLen - i))
		s = append(s, bit != 0)
	}

	parsed, err := ParseDaysOfWeek(s)
	if err != nil {
		return err
	}
	*w = *parsed

	return nil
}

func (src DaysOfWeek) BitsValue() (pgtype.Bits, error) {
	var acc uint8
	bitsLen := 7
	for i, b := range src.DaysOfWeek() {
		if b {
			acc += 2 << (bitsLen - (i + 1))
		}
	}
	return pgtype.Bits{
		Bytes: []byte{acc},
		Len:   int32(bitsLen),
		Valid: true,
	}, nil
}

var _ pgtype.BitsScanner = &DaysOfWeek{}
var _ pgtype.BitsValuer = DaysOfWeek{}

type NullDaysOfWeek struct {
	DaysOfWeek DaysOfWeek
	Valid      bool
}

func (nd *NullDaysOfWeek) ScanBits(bits pgtype.Bits) error {
	if !bits.Valid {
		nd.DaysOfWeek, nd.Valid = DaysOfWeek{}, false
		return nil
	}

	err := nd.DaysOfWeek.ScanBits(bits)
	if err != nil {
		nd.Valid = false
		return err
	}

	nd.Valid = true
	return nil
}

func (nd NullDaysOfWeek) BitsValue() (pgtype.Bits, error) {
	if !nd.Valid {
		return pgtype.Bits{}, nil
	}
	return nd.DaysOfWeek.BitsValue()
}

var _ pgtype.BitsScanner = &NullDaysOfWeek{}
var _ pgtype.BitsValuer = NullDaysOfWeek{}
