package double_test

import (
	"testing"

	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elim/double"
)

func TestDoubleElimination(t *testing.T) {
	d := tourney.Division{}
	d.MakeTeams(50)
	e := double.Elimination{}

	game := e.Generate(d)

	if game == nil {
		t.Error("Did not return game")
	}
}

func BenchmarkDoubleElimination(b *testing.B) {
	d := tourney.Division{}
	d.MakeTeams(b.N)

	b.ResetTimer()
	e := double.Elimination{}

	e.Generate(d)
}
