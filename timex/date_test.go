package timex

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		str  string
		want Date // if zero, expect an error
	}{
		{"1962-11-02", Date{1962, 11, 2}},
		{"1962-12-31", Date{1962, 12, 31}},
		{"0003-02-04", Date{3, 2, 4}},
		{"999-01-26", Date{}},
		{"", Date{}},
		{"1962-01-02x", Date{}},
	}

	for _, tt := range tests {
		got, err := ParseDate(tt.str)
		assert.Equal(t, tt.want, got)
		if got == (Date{}) {
			assert.NotNil(t, err)
		}
	}
}

func TestDateBefore(t *testing.T) {
	tests := []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{1962, 12, 31}, Date{1963, 1, 1}, true},
		{Date{1962, 1, 1}, Date{1962, 2, 1}, true},
		{Date{1962, 1, 1}, Date{1962, 1, 1}, false},
		{Date{1962, 12, 30}, Date{1962, 12, 31}, true},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.d1.Before(tt.d2))
	}
}

func TestDateAfter(t *testing.T) {
	tests := []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{1962, 12, 31}, Date{1963, 1, 1}, false},
		{Date{1962, 1, 1}, Date{1962, 2, 1}, false},
		{Date{1962, 1, 1}, Date{1962, 1, 1}, false},
		{Date{1962, 12, 30}, Date{1962, 12, 31}, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.d1.After(tt.d2))
	}
}

func TestDateCompare(t *testing.T) {
	tests := []struct {
		d1, d2 Date
		want   int
	}{
		{Date{1962, 12, 31}, Date{1963, 1, 1}, -1},
		{Date{1962, 1, 1}, Date{1962, 1, 1}, 0},
		{Date{1962, 12, 31}, Date{1962, 12, 30}, +1},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.d1.Compare(tt.d2))
	}
}

func TestMarshalJSON(t *testing.T) {
	got, err := json.Marshal(Date{2023, 5, 4})
	assert.Nil(t, err)
	assert.Equal(t, `"2023-05-04"`, string(got))
}

func TestUnmarshalJSON(t *testing.T) {
	var d Date
	err := json.Unmarshal([]byte(`"2023-05-04"`), &d)
	assert.Nil(t, err)
	assert.Equal(t, Date{2023, 5, 4}, d)
}

func TestMarshalXMLAttr(t *testing.T) {
	type Foo struct {
		Bar Date `xml:"bar,attr"`
	}
	got, err := xml.Marshal(Foo{Date{2023, 5, 4}})
	assert.Nil(t, err)
	assert.Equal(t, `<Foo bar="2023-05-04"></Foo>`, string(got))
}

func TestUnmarshalXML(t *testing.T) {
	var d Date
	err := xml.Unmarshal([]byte(`<Date>2023-05-04</Date>`), &d)
	assert.Nil(t, err)
	assert.Equal(t, Date{2023, 5, 4}, d)
}
