package robin

import "github.com/michael4d45/tourney"

type Game struct {
	Team1 tourney.Team
	Team2 tourney.Team

	Order int
	Round int
}