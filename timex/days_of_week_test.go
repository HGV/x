package timex

import (
	"testing"

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
