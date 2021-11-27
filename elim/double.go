package elim

import (
	"math"
	"sort"

	"github.com/michael4d45/tourney"
)

type DoubleGame struct {
	team1 *tourney.Team
	team2 *tourney.Team
	round int

	order int

	nextWinGame  *DoubleGame
	nextLoseGame *DoubleGame

	prevGame1 *DoubleGame
	prevGame2 *DoubleGame

	bracket string
}

func (g *DoubleGame) SetTeam1(team *tourney.Team) {
	g.team1 = team
}

func (g *DoubleGame) SetTeam2(team *tourney.Team) {
	g.team2 = team
}

func (g *DoubleGame) Teams() (*tourney.Team, *tourney.Team) {
	return g.team1, g.team2
}

func (g *DoubleGame) SetRound(round int) {
	g.round = round
}

func (g *DoubleGame) SetOrder(order int) {
	g.order = order
}

func (g *DoubleGame) Order() int {
	return g.order
}

func (g *DoubleGame) Round() int {
	return g.round
}

func (g *DoubleGame) Bracket() string {
	return g.bracket
}

func (g *DoubleGame) prevGame(game Game) *DoubleGame {
	team1, team2 := game.Teams()
	dGame := &DoubleGame{
		team1: team1,
		team2: team2,
		round: game.Round(),
		nextWinGame: g,
	}
	return dGame
}

func (g *DoubleGame) SetPrevGame1(game Game) {
	g.prevGame1 = g.prevGame(game)
}

func (g *DoubleGame) SetPrevGame2(game Game) {
	g.prevGame2 = g.prevGame(game)
}

func (g *DoubleGame) seed(e *generator, team1 *tourney.Team, team2 *tourney.Team, round int, oppositeRound float64) *DoubleGame {
	team1, team2 = e.teamOrder(team1, team2)
	game := &DoubleGame{
		team1: team1,
		team2: team2,
		round: round,
		nextWinGame: g,
	}

	e.seed(game, oppositeRound+1, round-1)

	return game
}

func (g *DoubleGame) Seed1(e *generator, team1 *tourney.Team, team2 *tourney.Team, round int, oppositeRound float64) {
	game := g.seed(e, team1, team2, round, oppositeRound)

	g.team1 = nil
	g.prevGame1 = game
}

func (g *DoubleGame) Seed2(e *generator, team1 *tourney.Team, team2 *tourney.Team, round int, oppositeRound float64) {
	game := g.seed(e, team1, team2, round, oppositeRound)

	g.team2 = nil
	g.prevGame2 = game
}

func (g *DoubleGame) PrevGame1() *Game {
	var game Game = g.prevGame1
	return &game
}

func (g *DoubleGame) PrevGame2() *Game {
	var game Game = g.prevGame2
	return &game
}

func (g *DoubleGame) NewGame() Game {
	return &DoubleGame{}
}

func GenerateDouble(d tourney.Division) *DoubleGame {
	g := &DoubleGame{}

	Generate(d, g)

	e := loserGen{}
	finalGame := e.GenerateLoser(d, g)

	return finalGame
}

type loserGen struct {
	rounds    [][]*DoubleGame
	roundsPos []int

	startGames    []*DoubleGame
	startGamesPos int

	order int
}

// loserBracketRoundCount calculates what round the last game is at.
func loserBracketRoundCount(n int) int {
	f := float64(n - 1)
	o := math.Log2(f)
	p := math.Log2(f * (8.0 / 3.0))
	return int(o) + int(p)
}

// shouldNumberTwice calculates if the looser bracket should be
// numbered twice or just once before starting the numbering cycle.
func shouldNumberTwice(count int) bool {
	n := math.Floor(math.Log2(float64(count)))
	o := math.Pow(2, n)
	p := o - (o / 4)
	c := float64(count)
	return (c <= o) && (c > p)
}

// Generate the
func (e *loserGen) GenerateLoser(division tourney.Division, game *DoubleGame) *DoubleGame {
	count := len(division.Teams)
	lRound := loserBracketRoundCount(count)

	// Last game if winner of wlGame lost.
	ifFirstLost := &DoubleGame{
		bracket: "WW",
		round:   lRound,
	}
	lRound--

	// Last game if winner of wlGame won.
	wlGame := &DoubleGame{
		bracket:      "WL",
		round:        lRound,
		nextWinGame:  ifFirstLost,
		nextLoseGame: ifFirstLost,
	}
	lRound--

	ifFirstLost.prevGame1 = wlGame
	ifFirstLost.prevGame2 = wlGame

	// Turn single elimination in double elimination games.
	gamesMap := map[*DoubleGame]struct{}{}
	game.doubleMap(&gamesMap)

	game.nextWinGame = wlGame

	wlGame.prevGame1 = game

	wRound := game.round
	e.rounds = make([][]*DoubleGame, wRound)
	e.roundsPos = make([]int, wRound)

	for gg := range gamesMap {
		e.rounds[gg.round-1] = append(e.rounds[gg.round-1], gg)
	}

	if wRound > 1 {
		// Collect starting games to hook together the winner's
		// bracket to the looser's bracket.
		sort.Slice(e.rounds[1], func(j, k int) bool {
			return e.rounds[1][j].order < e.rounds[1][k].order
		})
		e.startGames = make([]*DoubleGame, len(e.rounds[1]))
		e.setStartGames()

		// Sort the games in order to get a deterministic
		// winner to loser connection.
		for i, rr := range e.rounds {
			lowToHigh := (i % 2) != 0
			sort.Slice(rr, func(j, k int) bool {
				if lowToHigh {
					return rr[k].order < rr[j].order
				} else {
					return rr[j].order < rr[k].order
				}
			})
		}

		// Generate the loosers bracket.
		e.loserBracket(wlGame, 2, lRound, wRound)
	}

	// Set numbers on the games.
	lGame := wlGame.prevGame2

	numberTwice := shouldNumberTwice(count)

	for gg := range gamesMap {
		gg.bracket = "W"
		gg.order = 0
	}

	e.numberGames(game, lGame, numberTwice)
	wlGame.order = e.order
	e.order++
	ifFirstLost.order = e.order
	e.order++

	return ifFirstLost
}

func (game *DoubleGame) doubleMap(gamesMap *map[*DoubleGame]struct{}) {
	if game == nil {
		return
	}
	(*gamesMap)[game] = struct{}{}
	game.prevGame1.doubleMap(gamesMap)
	game.prevGame2.doubleMap(gamesMap)
}

// loserBracket generates the losers bracket.
func (e *loserGen) loserBracket(nextWinGame *DoubleGame, prevGame int, lRound int, wRound int) {
	// Use different rules when on the first round.
	if wRound == 2 {
		e.firstRoundGames(nextWinGame, prevGame, lRound)
		return
	}

	// Get a game from the winners bracket.
	game := e.takeFirstGame(wRound)

	// Create a loser bracket game that connects to
	// the winners bracket.
	fromWinnerGame := &DoubleGame{
		round:       lRound,
		nextWinGame: nextWinGame,
		prevGame1:   game,
		bracket:     "L",
	}
	lRound--

	if prevGame == 1 {
		nextWinGame.prevGame1 = fromWinnerGame
	}
	if prevGame == 2 {
		nextWinGame.prevGame2 = fromWinnerGame
	}

	losersPlayGame := &DoubleGame{
		round:       lRound,
		nextWinGame: fromWinnerGame,
		bracket:     "L",
	}
	lRound--

	fromWinnerGame.prevGame2 = losersPlayGame

	e.loserBracket(losersPlayGame, 1, lRound, wRound-1)
	e.loserBracket(losersPlayGame, 2, lRound, wRound-1)
}

func (e *loserGen) firstRoundGames(nextWinGame *DoubleGame, prevGame int, lRound int) {
	prevGame1 := e.takeFirststartGames()
	prevGame2 := e.takeFirstGame(2)
	game := prevGame2

	if prevGame1.round != 2 {
		game = &DoubleGame{
			round:       lRound,
			nextWinGame: nextWinGame,
			prevGame1:   prevGame1,
			prevGame2:   prevGame2,
			bracket:     "L",
		}

		if prevGame1.bracket == "L" {
			prevGame1.nextWinGame = game
		} else {
			prevGame1.nextLoseGame = game
		}
		if prevGame2.bracket == "L" {
			prevGame2.nextWinGame = game
		} else {
			prevGame2.nextLoseGame = game
		}
	}

	if prevGame == 1 {
		nextWinGame.prevGame1 = game
	}
	if prevGame == 2 {
		nextWinGame.prevGame2 = game
	}
}

func (e *loserGen) setStartGames() {
	for i, game := range e.rounds[1] {
		if game.prevGame1 != nil && game.prevGame2 == nil {
			e.startGames[i] = game.prevGame1
		} else if game.prevGame1 == nil && game.prevGame2 != nil {
			e.startGames[i] = game.prevGame2
		} else if game.prevGame1 != nil && game.prevGame2 != nil {
			game1 := game.prevGame1
			game2 := game.prevGame2
			game3 := &DoubleGame{
				round:     1,
				bracket:   "L",
				prevGame1: game1,
				prevGame2: game2,
			}
			e.startGames[i] = game3
			game1.nextLoseGame = game3
			game2.nextLoseGame = game3
		} else {
			game.bracket = "swap"
			e.startGames[i] = game
		}
	}
	jStart := len(e.startGames) - 1
	for i, game1 := range e.startGames {
		if game1.bracket == "swap" {
			game1.bracket = "W"

			for j := jStart; j > i; j-- {
				game2 := e.startGames[j]
				if game2.bracket == "swap" {
					game2.bracket = "W"

					e.startGames[i] = game2
					e.startGames[j] = game1
					jStart = j
					break
				}
			}
		}
	}
}

func (e *loserGen) numberGames(game *DoubleGame, lGame *DoubleGame, numberTwice bool) {
	wRound := 1
	lRound := 1

	e.order = 1

	e.numberGame(game, wRound, "W")
	wRound++

	e.numberGame(game, wRound, "W")
	wRound++

	e.numberGame(lGame, lRound, "L")
	lRound++

	if numberTwice {
		e.numberGame(lGame, lRound, "L")
		lRound++
	}

	for {
		checkFringe := e.numberGame(game, wRound, "W")
		wRound++

		if checkFringe != "not fringe" {
			break
		}

		e.numberGame(lGame, lRound, "L")
		lRound++

		e.numberGame(lGame, lRound, "L")
		lRound++
	}

	for {
		checkFringe := e.numberGame(lGame, lRound, "L")
		lRound++
		if checkFringe != "not fringe" {
			break
		}
	}
}

func (e *loserGen) numberGame(game *DoubleGame, round int, bracket string) string {
	if game == nil || game.bracket != bracket || game.order != 0 {
		return "null"
	}
	game1 := e.numberGame(game.prevGame1, round, bracket)
	game2 := e.numberGame(game.prevGame2, round, bracket)
	if game1 == "null" && game2 == "null" && game.round == round {
		game.order = e.order
		e.order++
		return "fringe"
	}
	return "not fringe"
}

func (e *loserGen) takeFirststartGames() *DoubleGame {
	game := e.startGames[e.startGamesPos]
	e.startGamesPos++
	return game
}

func (e *loserGen) takeFirstGame(round int) *DoubleGame {
	game := e.rounds[round-1][e.roundsPos[round-1]]
	e.roundsPos[round-1]++
	return game
}
