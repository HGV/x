package slicesx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartition(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }

	input := []int{1, 2, 3, 4, 5, 6}
	wantEven := []int{2, 4, 6}
	wantOdd := []int{1, 3, 5}

	even, odd := Partition(input, isEven)

	assert.EqualValues(t, wantEven, even)
	assert.EqualValues(t, wantOdd, odd)
}

func TestPartition_Empty(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }

	input := []int{}
	even, odd := Partition[[]int](input, isEven)

	assert.Len(t, even, 0)
	assert.Len(t, odd, 0)
}

func TestPartition_Nil(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }

	even, odd := Partition[[]int](nil, isEven)

	assert.Len(t, even, 0)
	assert.Len(t, odd, 0)
}
