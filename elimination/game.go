package elimination

import (
	"github.com/michael4d45/tourney"
)

type Game struct {
	Team1 *tourney.Team
	Team2 *tourney.Team
	Round int

	PrevGame1 *Game
	PrevGame2 *Game
}