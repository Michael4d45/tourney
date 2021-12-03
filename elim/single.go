package elim

import (
	"math"

	"github.com/michael4d45/tourney"
)

type singleGenerator struct {
	teamMap  map[int]*tourney.Team
	topOrder string
	order    int

	rounds [][]*Game
}

func roundCount(len int) int {
	p := math.Floor(math.Log2(float64(len)))
	highestPower := math.Pow(2, p)
	if int(highestPower) == len {
		return int(p)
	}
	return int(p + 1)
}

func (s *singleGenerator) generateSingle(division tourney.Division, topOrder string) *Game {
	if len(division.Teams) <= 1 {
		return nil
	}

	s.topOrder = topOrder

	s.teamMap = map[int]*tourney.Team{}

	for _, tt := range division.Teams {
		s.teamMap[tt.Seed] = tt
	}

	numRounds := roundCount(len(s.teamMap))

	s.rounds = make([][]*Game, numRounds)

	team1 := s.teamMap[1]
	team2 := s.otherTeam(team1, 3)
	game := s.newGame(team1, team2, numRounds, nil)

	s.seed(game, 2, numRounds-1)

	s.numberGames(game)

	return game
}

func (s *singleGenerator) seed(winGame *Game, oppositeRound float64, round int) {
	if round < 1 {
		return
	}
	roundSeed := int(math.Pow(2, oppositeRound)) + 1

	team3 := s.otherTeam(winGame.Team1, roundSeed)
	if team3 != nil {
		newGame := s.newGame(winGame.Team1, team3, round, winGame)
		winGame.Team1 = nil
		winGame.PrevGame1 = newGame
		s.seed(newGame, oppositeRound+1, round-1)
	}

	team4 := s.otherTeam(winGame.Team2, roundSeed)
	if team4 != nil {
		newGame := s.newGame(team4, winGame.Team2, round, winGame)
		winGame.Team2 = nil
		winGame.PrevGame2 = newGame
		s.seed(newGame, oppositeRound+1, round-1)
	}
}

func (s *singleGenerator) newGame(team1, team2 *tourney.Team, round int, nextGame *Game) *Game {
	game := &Game{
		Team1:       team1,
		Team2:       team2,
		Round:       round,
		NextWinGame: nextGame,
	}
	s.rounds[round-1] = append(s.rounds[round-1], game)
	s.orderTeams(game)
	return game
}

func (s *singleGenerator) otherTeam(team1 *tourney.Team, roundSeed int) *tourney.Team {
	seed := roundSeed - team1.Seed
	return s.teamMap[seed]
}

func (s *singleGenerator) orderTeams(Game *Game) {
	if Game.Team1 == nil || Game.Team2 == nil {
		return
	}
	var even, odd, high, low *tourney.Team
	if Game.Team1.Seed%2 == 0 {
		even = Game.Team1
		odd = Game.Team2
	} else {
		even = Game.Team2
		odd = Game.Team1
	}
	if Game.Team1.Seed > Game.Team2.Seed {
		high = Game.Team1
		low = Game.Team2
	} else {
		high = Game.Team2
		low = Game.Team1
	}
	switch s.topOrder {
	case "even":
		Game.Team1 = even
		Game.Team2 = odd
	case "odd":
		Game.Team1 = odd
		Game.Team2 = even
	case "high":
		Game.Team1 = high
		Game.Team2 = low
	case "low":
		Game.Team1 = low
		Game.Team2 = high
	}
}

func (s *singleGenerator) numberGames(Game *Game) {
	round := 1
	s.order = 1
	for ; s.numberGame(Game, round) == "not fringe"; round++ {
	}
}

func (s *singleGenerator) numberGame(Game *Game, round int) string {
	if Game == nil || Game.Order != 0 {
		return "null"
	}
	Game1 := s.numberGame(Game.PrevGame1, round)
	Game2 := s.numberGame(Game.PrevGame2, round)
	if Game1 == "null" && Game2 == "null" && Game.Round == round {
		Game.Order = s.order
		s.order++
		return "fringe"
	}

	return "not fringe"
}
