package single

import (
	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elimination"
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

func copy(g *elimination.Game, games map[*elimination.Game]*Game, prevGame *Game) (*Game, map[*elimination.Game]*Game) {
	if g == nil {
		return nil, games
	}
	game, exists := games[g]
	if exists {
		return game, games
	}
	game = &Game{
		Team1:    g.Team1,
		Team2:    g.Team2,
		Round:    g.Round,
		NextGame: prevGame,
	}
	games[g] = game
	game.PrevGame1, games = copy(g.PrevGame1, games, game)
	game.PrevGame2, games = copy(g.PrevGame2, games, game)

	return game, games
}
