package elim

import (
	"github.com/michael4d45/tourney"
)

type Game struct {
	Team1 *tourney.Team
	Team2 *tourney.Team
	Round int

	Order int

	NextWinGame  *Game
	NextLoseGame *Game

	PrevGame1 *Game
	PrevGame2 *Game

	Bracket string
}
