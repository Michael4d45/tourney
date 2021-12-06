package strings

import (
	"strconv"

	"github.com/michael4d45/tourney/robin"
)

func Robin(games robin.Games) string {
	s := ""
	for _, rr := range games.Rounds {
		for _, gg := range rr {
			s += strconv.Itoa(gg.Team1.Seed)
			if gg.Team1.Seed < 10 && gg.Team2.Seed >= 10 {
				s += " "
			} 
			s += " | "
		}
		s += "\n"
		for _, gg := range rr {
			s += strconv.Itoa(gg.Team2.Seed)
			if gg.Team2.Seed < 10 && gg.Team1.Seed >= 10 {
				s += " "
			} 
			s += " | "
		}
		s += "\n\n"
	}

	return s
}
