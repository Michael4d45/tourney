package brackets

import (
	"math"
)

type Elimination struct {
	teamMap map[int]*Team
	topOrder string
}

func winnerBracketRoundCount(len int) int {
	p := math.Floor(math.Log2(float64(len)))
	highestPower := math.Pow(2, p)
	if int(highestPower) == len {
		return int(p)
	}
	return int(p + 1)
}

func (e *Elimination) Generate(division Division) *Game {
	if len(division.Teams) <= 1 {
		return nil
	}

	e.topOrder = "even" // [even/odd],[lower/higher],random

	e.teamMap = map[int]*Team{}

	for _, team := range division.Teams {
		e.teamMap[team.Seed] = team
	}

	numRounds := winnerBracketRoundCount(len(e.teamMap))

	team1 := e.teamMap[1]
	team2 := e.otherTeam(team1, 3)
	game := &Game{
		team1: team1,
		team2: team2,
		round: numRounds,
	}

	e.seed(game, 2, numRounds-1)


	return game
}

func (e *Elimination) seed(winGame *Game, oppositeRound float64, round int) {
	if round <= 1 {
		return
	}
	roundSeed := int(math.Pow(2, oppositeRound)) + 1

	team4 := e.otherTeam(winGame.team1, roundSeed)
	if team4 != nil {
		game1 := &Game{
			team1:         winGame.team1,
			team2:         team4,
			round:         round,
			nextWinGame: winGame,
		}
		winGame.team1 = nil
		e.seed(game1, oppositeRound+1, round-1)
		winGame.prevGame1 = game1
	}

	team3 := e.otherTeam(winGame.team2, roundSeed)
	if team3 != nil {
		game2 := &Game{
			team1:         team3,
			team2:         winGame.team2,
			round:         round,
			nextWinGame: winGame,
		}
		winGame.team2 = nil
		e.seed(game2, oppositeRound+1, round-1)
		winGame.prevGame2 = game2
	}
}

func (e *Elimination) otherTeam(team1 *Team, roundSeed int) *Team {
	seed := roundSeed - team1.Seed
	return e.teamMap[seed]
}

func (e *Elimination) order(game *Game) {
	switch(e.topOrder) {
	case "even":
		
	}
}
