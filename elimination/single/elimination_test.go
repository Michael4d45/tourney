package single_test

import (
	"testing"

	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elimination/single"
)

func TestElimination(t *testing.T) {
	d := tourney.Division{}
	d.MakeTeams(50)
	e := single.Elimination{}

	game := e.Generate(d)

	if game == nil {
		t.Error("Did not return game")
	}
}

func BenchmarkElimination(b *testing.B) {
	d := tourney.Division{}
	d.MakeTeams(b.N)

	b.ResetTimer()
	e := single.Elimination{}

	e.Generate(d)
}
