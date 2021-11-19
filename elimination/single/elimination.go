package single

import (
	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elimination"
)

type Elimination struct {
	gameNum int
}

func (e *Elimination) Generate(division tourney.Division) *Game {
	elim := elimination.Elimination{}
	eGame := elim.Generate(division)
	if eGame == nil {
		return nil
	}

	wGame, _ := copy(eGame, map[*elimination.Game]*Game{}, nil)

	e.numberGames(wGame)

	return wGame
}

func (e *Elimination) numberGames(game *Game) {
	round := 1
	e.gameNum = 1
	for ;e.numberGame(game, round) == "not fringe";round++ {}
}

func (e *Elimination) numberGame(game *Game, round int) string {
	if (game == nil || game.GameNum != 0) {
		return "null"
	}
	game1 := e.numberGame(game.PrevGame1, round)
	game2 := e.numberGame(game.PrevGame2, round)
	if game1 == "null" && game2 == "null" && game.Round == round {
		game.GameNum = e.gameNum
		e.gameNum++
		return "fringe"
	}

	return "not fringe"
}