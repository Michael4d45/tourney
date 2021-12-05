package elim_test

import (
	"testing"

	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elim"
)

func TestTopOrder(t *testing.T) {
	d := tourney.Division{}
	d.MakeTeams(5)

	orders := []string{"odd", "even", "high", "low"}

	for _, topOrder := range orders {
		games := elim.Generate(d, topOrder, "double")
		if games.FinalGame == nil {
			t.Error("Did not return game for ", topOrder)
		}
	}
}

func TestWorkingTeams(t *testing.T) {
	topOrder := "odd"

	teams := []int{2, 3, 4, 100}

	for _, teamsCount := range teams {
		d := tourney.Division{}
		d.MakeTeams(teamsCount)

		games := elim.Generate(d, topOrder, "double")
		if games.FinalGame == nil {
			t.Error("Did not return game for ", teamsCount)
		}
	}
}

func TestFailingTeams(t *testing.T) {
	topOrder := "odd"

	teams := []int{0, 1}

	for _, teamsCount := range teams {
		d := tourney.Division{}
		d.MakeTeams(teamsCount)

		games := elim.Generate(d, topOrder, "double")
		if games.FinalGame == nil {
			t.Error("Did return game for ", teamsCount)
		}
	}
}
