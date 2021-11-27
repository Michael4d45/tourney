package elim

import (
	"math"

	"github.com/michael4d45/tourney"
)

type Game interface {
	SetTeam1(*tourney.Team)
	SetTeam2(*tourney.Team)

	Teams() (*tourney.Team, *tourney.Team)

	SetRound(int)
	SetOrder(int)

	Order() int
	Round() int

	SetPrevGame1(Game)
	SetPrevGame2(Game)

	Seed1(*generator, *tourney.Team, *tourney.Team, int, float64)
	Seed2(*generator, *tourney.Team, *tourney.Team, int, float64)

	PrevGame1() *Game
	PrevGame2() *Game

	NewGame() Game
}

type generator struct {
	teamMap  map[int]*tourney.Team
	topOrder string
	Order    int
}

func roundCount(len int) int {
	p := math.Floor(math.Log2(float64(len)))
	highestPower := math.Pow(2, p)
	if int(highestPower) == len {
		return int(p)
	}
	return int(p + 1)
}

func Generate(d tourney.Division, g Game) {
	e := generator{}
	e.topOrder = "odd" // [even/odd],[lower/higher],random

	e.teamMap = make(map[int]*tourney.Team)
	for _, tt := range d.Teams {
		e.teamMap[tt.Seed] = tt
	}

	numRounds := roundCount(len(e.teamMap))

	team1 := e.teamMap[1]
	team2 := e.otherTeam(team1, 3)

	g.SetRound(numRounds)

	team1, team2 = e.teamOrder(team1, team2)
	g.SetTeam1(team1)
	g.SetTeam2(team2)

	e.seed(g, 2, numRounds-1)

	e.numberGames(&g)
}

func (e *generator) seed(winGame Game, oppositeRound float64, round int) {
	if round < 1 {
		return
	}
	roundSeed := int(math.Pow(2, oppositeRound)) + 1

	team1, team2 := winGame.Teams()

	team4 := e.otherTeam(team1, roundSeed)
	if team4 != nil {
		winGame.Seed1(e, team1, team4, round, oppositeRound)
	}

	team3 := e.otherTeam(team2, roundSeed)
	if team3 != nil {
		winGame.Seed2(e, team3, team2, round, oppositeRound)
	}
}

func (e *generator) otherTeam(team1 *tourney.Team, roundSeed int) *tourney.Team {
	seed := roundSeed - team1.Seed
	return e.teamMap[seed]
}

func (e *generator) teamOrder(team1 *tourney.Team, team2 *tourney.Team) (*tourney.Team, *tourney.Team) {
	if team1 == nil || team2 == nil {
		return team1, team2
	}
	var even, odd, high, low *tourney.Team
	if team1.Seed%2 == 0 {
		even = team1
		odd = team2
	} else {
		even = team2
		odd = team1
	}
	if team1.Seed > team2.Seed {
		high = team1
		low = team2
	} else {
		high = team2
		low = team1
	}
	switch e.topOrder {
	case "even":
		return even, odd
	case "odd":
		return odd, even
	case "high":
		return high, low
	case "low":
		return low, high
	default:
		return team1, team2
	}
}

func (e *generator) numberGames(game *Game) {
	round := 1
	e.Order = 1
	for {
		checkFringe := e.numberGame(game, round)
		round++
		if checkFringe != "not fringe" {
			break
		}
	}
}

func IsElimNil(i Game) bool {
	var ret bool
	switch i.(type) {
	case *DoubleGame:
		v := i.(*DoubleGame)
		ret = v == nil
	case *SingleGame:
		v := i.(*SingleGame)
		ret = v == nil
	}
	return ret
}

func (e *generator) numberGame(game *Game, round int) string {
	if IsElimNil(*game) || (*game).Order() != 0 {
		return "null"
	}
	game1 := e.numberGame((*game).PrevGame1(), round)
	game2 := e.numberGame((*game).PrevGame2(), round)
	if game1 == "null" && game2 == "null" && (*game).Round() == round {
		(*game).SetOrder(e.Order)
		e.Order++
		return "fringe"
	}

	return "not fringe"
}
