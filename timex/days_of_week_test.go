package timex

import (
	"testing"
	"testing/quick"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestParseDaysOfWeek(t *testing.T) {
	valid := []bool{true, false, true, false, true, false, true}
	dow, err := ParseDaysOfWeek(valid)
	assert.NoError(t, err)
	assert.Equal(t, valid, dow.DaysOfWeek())

	invalid := []bool{true, false}
	_, err = ParseDaysOfWeek(invalid)
	assert.Error(t, err)
}

func TestHas_PropertyBased(t *testing.T) {
	f := func(mo, tu, we, th, fr, sa, su bool) bool {
		dow := DaysOfWeek{Mo: mo, Tu: tu, We: we, Th: th, Fr: fr, Sa: sa, Su: su}
		want := []bool{su, mo, tu, we, th, fr, sa}
		for i, expected := range want {
			actual := dow.Has(time.Weekday(i))
			if actual != expected {
				return false
			}
		}
		return true
	}
	assert.NoError(t, quick.Check(f, nil))
}

func TestHas_InvalidWeekday(t *testing.T) {
	var dow DaysOfWeek
	assert.False(t, dow.Has(time.Weekday(8)))
}

func TestHas_ValidWeekdays(t *testing.T) {
	tests := []struct {
		name string
		dow  DaysOfWeek
		want []bool
	}{
		{
			name: "all false",
			dow:  DaysOfWeek{},
			want: []bool{false, false, false, false, false, false, false},
		},
		{
			name: "all true",
			dow:  DaysOfWeek{Mo: true, Tu: true, We: true, Th: true, Fr: true, Sa: true, Su: true},
			want: []bool{true, true, true, true, true, true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, expected := range tt.want {
				actual := tt.dow.Has(time.Weekday(i))
				assert.Equal(t, actual, expected)
			}
		})
	}
}

func TestDaysOfWeek_BitsConversion(t *testing.T) {
	original := DaysOfWeek{Mo: true, We: true, Fr: true, Su: true}

	// Convert to bits
	bits, err := original.BitsValue()
	assert.NoError(t, err)
	assert.True(t, bits.Valid)
	assert.Equal(t, int32(7), bits.Len)

	// Read from bits
	var decoded DaysOfWeek
	err = decoded.ScanBits(bits)
	assert.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestDaysOfWeek_Compare(t *testing.T) {
	a := DaysOfWeek{Mo: true, Tu: false}
	b := DaysOfWeek{Mo: true, Tu: true}
	c := DaysOfWeek{Mo: true, Tu: false}

	assert.NotZero(t, a.Compare(b))
	assert.Equal(t, 0, a.Compare(c))
}

func TestNullDaysOfWeek_BitsConversion(t *testing.T) {
	original := NullDaysOfWeek{
		DaysOfWeek: DaysOfWeek{Mo: true, Sa: true},
		Valid:      true,
	}

	// Convert to bits
	bits, err := original.BitsValue()
	assert.NoError(t, err)

	// Decode from bits
	var decoded NullDaysOfWeek
	err = decoded.ScanBits(bits)
	assert.NoError(t, err)

	assert.True(t, decoded.Valid)
	assert.Equal(t, original.DaysOfWeek, decoded.DaysOfWeek)
}

func TestNullDaysOfWeek_Invalid(t *testing.T) {
	invalidBits := pgtype.Bits{Valid: false}

	var nd NullDaysOfWeek
	err := nd.ScanBits(invalidBits)
	assert.NoError(t, err)
	assert.False(t, nd.Valid)

	bits, err := nd.BitsValue()
	assert.NoError(t, err)
	assert.False(t, bits.Valid)
}
