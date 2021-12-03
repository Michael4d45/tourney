package elim_test

import (
	"testing"

	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elim"
)

func TestDouble(t *testing.T) {
	d := tourney.Division{}
	d.MakeTeams(5)

	topOrder := "odd"
	elimType := "double"

	game := elim.Generate(d, topOrder, elimType)
	if game == nil {
		t.Error("Did not return game for double")
	}
}