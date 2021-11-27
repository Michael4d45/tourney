package elim

import (
	"github.com/michael4d45/tourney"
)

type SingleGame struct {
	team1 *tourney.Team
	team2 *tourney.Team
	round int

	order int

	nextGame  *SingleGame

	prevGame1 *SingleGame
	prevGame2 *SingleGame

	bracket string
}

func (g *SingleGame) SetTeam1(team *tourney.Team) {
	g.team1 = team
}

func (g *SingleGame) SetTeam2(team *tourney.Team) {
	g.team2 = team
}

func (g *SingleGame) Teams() (*tourney.Team, *tourney.Team) {
	return g.team1, g.team2
}

func (g *SingleGame) SetRound(round int) {
	g.round = round
}

func (g *SingleGame) SetOrder(order int) {
	g.order = order
}

func (g *SingleGame) Order() int {
	return g.order
}

func (g *SingleGame) Round() int {
	return g.round
}

func (g *SingleGame) Bracket() string {
	return g.bracket
}

func (g *SingleGame) prevGame(game Game) *SingleGame {
	team1, team2 := game.Teams()
	sGame := &SingleGame{
		team1: team1,
		team2: team2,
		round: game.Round(),
		nextGame: g,
	}
	return sGame
}

func (g *SingleGame) SetPrevGame1(game Game) {
	g.prevGame1 = g.prevGame(game)
}

func (g *SingleGame) SetPrevGame2(game Game) {
	g.prevGame2 = g.prevGame(game)
}

func (g *SingleGame) seed(e *generator, team1 *tourney.Team, team2 *tourney.Team, round int, oppositeRound float64) *SingleGame {
	team1, team2 = e.teamOrder(team1, team2)
	game := &SingleGame{
		team1: team1,
		team2: team2,
		round: round,
		nextGame: g,
	}

	e.seed(game, oppositeRound+1, round-1)

	return game
}

func (g *SingleGame) Seed1(e *generator, team1 *tourney.Team, team2 *tourney.Team, round int, oppositeRound float64) {
	game := g.seed(e, team1, team2, round, oppositeRound)

	g.team1 = nil
	g.prevGame1 = game
}

func (g *SingleGame) Seed2(e *generator, team1 *tourney.Team, team2 *tourney.Team, round int, oppositeRound float64) {
	game := g.seed(e, team1, team2, round, oppositeRound)

	g.team2 = nil
	g.prevGame2 = game
}

func (g *SingleGame) PrevGame1() *Game {
	var game Game = g.prevGame1
	return &game
}

func (g *SingleGame) PrevGame2() *Game {
	var game Game = g.prevGame2
	return &game
}

func (g *SingleGame) NewGame() Game {
	return &SingleGame{}
}

func GenerateSingle(d tourney.Division) *SingleGame {
	g := &SingleGame{}

	Generate(d, g)

	return g
}


