package elimination

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

	SetPrevGame1(*Game)
	SetPrevGame2(*Game)

	PrevGame1() *Game
	PrevGame2() *Game
}

type Generator struct {
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

func Generate(d tourney.Division, g *Game) {
	e := Generator{}
	e.topOrder = "odd" // [even/odd],[lower/higher],random

	e.teamMap = make(map[int]*tourney.Team)
	for _, tt := range d.Teams {
		e.teamMap[tt.Seed] = tt
	}

	numRounds := roundCount(len(e.teamMap))

	team1 := e.teamMap[1]
	team2 := e.otherTeam(team1, 3)

	(*g).SetTeam1(team1)
	(*g).SetTeam2(team2)
	(*g).SetRound(numRounds)

	e.order(g)

	e.seed(g, 2, numRounds-1)

	e.numberGames(g)
}
func (e *Generator) seed(winGame *Game, oppositeRound float64, round int) {
	if round < 1 {
		return
	}
	roundSeed := int(math.Pow(2, oppositeRound)) + 1

	team1, team2 := (*winGame).Teams()

	team4 := e.otherTeam(team1, roundSeed)
	if team4 != nil {
		var game1 Game
		game1.SetTeam1(team1)
		game1.SetTeam2(team4)
		game1.SetRound(round)

		e.order(&game1)
		(*winGame).SetTeam1(nil)
		e.seed(&game1, oppositeRound+1, round-1)
		(*winGame).SetPrevGame1(&game1)
	}

	team3 := e.otherTeam(team2, roundSeed)
	if team3 != nil {
		var game2 Game
		game2.SetTeam1(team3)
		game2.SetTeam2(team2)
		game2.SetRound(round)
		
		e.order(&game2)
		(*winGame).SetTeam2(nil)
		e.seed(&game2, oppositeRound+1, round-1)
		(*winGame).SetPrevGame2(&game2)
	}
}

func (e *Generator) otherTeam(team1 *tourney.Team, roundSeed int) *tourney.Team {
	seed := roundSeed - team1.Seed
	return e.teamMap[seed]
}

func (e *Generator) order(game *Game) {
	team1, team2 := (*game).Teams()
	if team1 == nil || team2 == nil {
		return
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
		(*game).SetTeam1(even)
		(*game).SetTeam2(odd)
	case "odd":
		(*game).SetTeam1(odd)
		(*game).SetTeam2(even)
	case "high":
		(*game).SetTeam1(high)
		(*game).SetTeam2(low)
	case "low":
		(*game).SetTeam1(low)
		(*game).SetTeam2(high)
	}
}

func (e *Generator) numberGames(game *Game) {
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

func (e *Generator) numberGame(game *Game, round int) string {
	if game == nil || (*game).Order() != 0 {
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
