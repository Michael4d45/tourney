package elimination

import (
	"math"

	"github.com/michael4d45/tourney"
)

type Elimination struct {
	teamMap  map[int]*tourney.Team
	topOrder string
}

func roundCount(len int) int {
	p := math.Floor(math.Log2(float64(len)))
	highestPower := math.Pow(2, p)
	if int(highestPower) == len {
		return int(p)
	}
	return int(p + 1)
}

func (e *Elimination) Generate(division tourney.Division) *Game {
	if len(division.Teams) <= 1 {
		return nil
	}

	e.topOrder = "odd" // [even/odd],[lower/higher],random

	e.teamMap = map[int]*tourney.Team{}

	for _, tt := range division.Teams {
		e.teamMap[tt.Seed] = tt
	}

	numRounds := roundCount(len(e.teamMap))

	team1 := e.teamMap[1]
	team2 := e.otherTeam(team1, 3)
	game := &Game{
		Team1: team1,
		Team2: team2,
		Round: numRounds,
	}
	e.order(game)

	e.seed(game, 2, numRounds-1)

	return game
}

func (e *Elimination) seed(winGame *Game, oppositeRound float64, round int) {
	if round < 1 {
		return
	}
	roundSeed := int(math.Pow(2, oppositeRound)) + 1

	team4 := e.otherTeam(winGame.Team1, roundSeed)
	if team4 != nil {
		game1 := &Game{
			Team1:    winGame.Team1,
			Team2:    team4,
			Round:    round,
		}
		e.order(game1)
		winGame.Team1 = nil
		e.seed(game1, oppositeRound+1, round-1)
		winGame.PrevGame1 = game1
	}

	team3 := e.otherTeam(winGame.Team2, roundSeed)
	if team3 != nil {
		game2 := &Game{
			Team1:    team3,
			Team2:    winGame.Team2,
			Round:    round,
		}
		e.order(game2)
		winGame.Team2 = nil
		e.seed(game2, oppositeRound+1, round-1)
		winGame.PrevGame2 = game2
	}
}

func (e *Elimination) otherTeam(team1 *tourney.Team, roundSeed int) *tourney.Team {
	seed := roundSeed - team1.Seed
	return e.teamMap[seed]
}

func (e *Elimination) order(game *Game) {
	if game.Team1 == nil || game.Team2 == nil {
		return
	}
	var even, odd, high, low *tourney.Team
	if game.Team1.Seed%2 == 0 {
		even = game.Team1
		odd = game.Team2
	} else {
		even = game.Team2
		odd = game.Team1
	}
	if game.Team1.Seed > game.Team2.Seed {
		high = game.Team1
		low = game.Team2
	} else {
		high = game.Team2
		low = game.Team1
	}
	switch e.topOrder {
	case "even":
		game.Team1 = even
		game.Team2 = odd
	case "odd":
		game.Team1 = odd
		game.Team2 = even
	case "high":
		game.Team1 = high
		game.Team2 = low
	case "low":
		game.Team1 = low
		game.Team2 = high
	}
}
