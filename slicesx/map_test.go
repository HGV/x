package slicesx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
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

	lowered := Map(nobelPrizes, func(n nobelPrize) string {
		return strings.ToLower(n.Laureate)
	})

	assert.EqualValues(t, []string{
		"luis leloir",
		"paul samuelson",
		"aleksandr solzhenitsyn",
		"norman borlaug",
		"hannes alfvén",
		"louis néel",
		"sir bernard katz",
		"ulf von euler",
		"julius axelrod",
	}, lowered)
}
