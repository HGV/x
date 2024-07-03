package slicesx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupByFunc(t *testing.T) {
	type nobelPrize struct {
		Category string
		Laureate string
	}
	nobelPrizes := []nobelPrize{
		{Category: "Chemistry", Laureate: "Luis Leloir"},
		{Category: "Economics", Laureate: "Paul Samuelson"},
		{Category: "Literature", Laureate: "Aleksandr Solzhenitsyn"},
		{Category: "Peace", Laureate: "Norman Borlaug"},
		{Category: "Physics", Laureate: "Hannes Alfvén"},
		{Category: "Physics", Laureate: "Louis Néel"},
		{Category: "Medicine", Laureate: "Sir Bernard Katz"},
		{Category: "Medicine", Laureate: "Ulf von Euler"},
		{Category: "Medicine", Laureate: "Julius Axelrod"},
	}

	byCategory := GroupByFunc(nobelPrizes, func(n nobelPrize) string {
		return n.Category
	})

	assert.Equal(t, 6, len(byCategory))
	assert.EqualValues(t, []nobelPrize{
		{Category: "Physics", Laureate: "Hannes Alfvén"},
		{Category: "Physics", Laureate: "Louis Néel"},
	}, byCategory["Physics"])
}
