package strings

import (
	"strconv"

	"github.com/michael4d45/tourney"
)

func Team(t tourney.Team) string {
	s := "Team: "
	s += "seed: " + strconv.Itoa(t.Seed)
	// s += ", d: " + strconv.Itoa(len(t.Division.Teams))
	return s
}

func Division(d tourney.Division) string {
	s := "division: "
	s += "\n"
	s += "\tteams: ["

	s += "\n"
	for i := 0; i < len(d.Teams); i++ {
		s += "\t\t"
		s += Team(*d.Teams[i])
		s += "\n"
	}
	s += "\t]"

	s += "\n"
	return s
}