package single

import (
	"github.com/michael4d45/tourney"
)

type Game struct {
	Team1 *tourney.Team
	Team2 *tourney.Team
	Round int

	GameNum int

	NextGame  *Game
	PrevGame1 *Game
	PrevGame2 *Game
}
