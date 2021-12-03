package elim

import (
	"github.com/michael4d45/tourney"
)

// [even/odd],[lower/higher], random
func Generate(division tourney.Division, topOrder string, elimType string) *Game {

	if elimType == "single" || elimType == "double" {
		s := singleGenerator{}
		g := s.generateSingle(division, topOrder)

		if elimType == "double" && g != nil {
			l := loserGen{}
			g = l.generateLoser(g, len(division.Teams), s.rounds)
		}

		return g
	}

	return nil
}
