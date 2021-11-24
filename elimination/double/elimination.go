package double

import (
	"math"
	"sort"

	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elimination/single"
)

type Elimination struct {
	rounds    [][]*Game
	roundsPos []int
	startGames    []*Game
	startGamesPos int
	order     int
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
func (e *Elimination) Generate(division tourney.Division) *Game {

	// Get the single elimination tournamnet.
	elim := single.Elimination{}
	eGame := elim.Generate(division)
	if eGame == nil {
		return nil
	}

	count := len(division.Teams)
	lRound := loserBracketRoundCount(count)

	// Last game if winner of wlGame lost.
	ifFirstLost := &Game{
		Bracket: "WW",
		Round:   lRound,
	}
	lRound--

	// Last game if winner of wlGame won.
	wlGame := &Game{
		Bracket:      "WL",
		Round:        lRound,
		NextWinGame:  ifFirstLost,
		NextLoseGame: ifFirstLost,
	}
	lRound--

	ifFirstLost.PrevGame1 = wlGame
	ifFirstLost.PrevGame2 = wlGame

	// Turn single elimination in double elimination games.
	wGame, gamesMap := copy(eGame, map[*single.Game]*Game{})

	wGame.NextWinGame = wlGame

	wlGame.PrevGame1 = wGame

	wRound := wGame.Round
	e.rounds = make([][]*Game, wRound)
	e.roundsPos = make([]int, wRound)

	for _, gg := range gamesMap {
		e.rounds[gg.Round-1] = append(e.rounds[gg.Round-1], gg)
	}

	if wRound > 1 {
		// Collect starting games to hook together the winner's
		// bracket to the looser's bracket.
		sort.Slice(e.rounds[1], func(j, k int) bool {
			return e.rounds[1][j].Order < e.rounds[1][k].Order
		})
		e.startGames = make([]*Game, len(e.rounds[1]))
		e.setStartGames()

		// Sort the games in order to get a deterministic
		// winner to loser connection.
		for i, rr := range e.rounds {
			lowToHigh := (i % 2) != 0
			sort.Slice(rr, func(j, k int) bool {
				if lowToHigh {
					return rr[k].Order < rr[j].Order
				} else {
					return rr[j].Order < rr[k].Order
				}
			})
		}

		// Actually generate the loosers bracket.
		e.loserBracket(wlGame, 2, lRound, wRound)
	}

	// Set numbers on the games.
	lGame := wlGame.PrevGame2

	numberTwice := shouldNumberTwice(count)

	for _, gg := range gamesMap {
		gg.Bracket = "W"
		gg.Order = 0
	}

	e.numberGames(wGame, lGame, numberTwice)
	wlGame.Order = e.order
	e.order++
	ifFirstLost.Order = e.order
	e.order++

	return ifFirstLost
}

// loserBracket generates the losers bracket. 
func (e *Elimination) loserBracket(nextWinGame *Game, prevGame int, lRound int, wRound int) {
	// Use different rules when on the first round.
	if wRound == 2 {
		e.firstRoundGames(nextWinGame, prevGame, lRound)
		return
	}

	// Get a game from the winners bracket.
	wGame := e.takeFirstGame(wRound)
	
	// Create a loser bracket game that connects to 
	// the winners bracket.
	fromWinnerGame := &Game{
		Round:       lRound,
		NextWinGame: nextWinGame,
		PrevGame1:   wGame,
		Bracket:     "L",
	}
	lRound--

	if prevGame == 1 {
		nextWinGame.PrevGame1 = fromWinnerGame
	}
	if prevGame == 2 {
		nextWinGame.PrevGame2 = fromWinnerGame
	}

	losersPlayGame := &Game{
		Round:       lRound,
		NextWinGame: fromWinnerGame,
		Bracket:     "L",
	}
	lRound--

	fromWinnerGame.PrevGame2 = losersPlayGame

	e.loserBracket(losersPlayGame, 1, lRound, wRound-1)
	e.loserBracket(losersPlayGame, 2, lRound, wRound-1)
}

func (e *Elimination) firstRoundGames(nextWinGame *Game, prevGame int, lRound int) {
	prevGame1 := e.takeFirststartGames()
	prevGame2 := e.takeFirstGame(2)
	game := prevGame2

	if prevGame1.Round != 2 {
		game = &Game{
			Round:       lRound,
			NextWinGame: nextWinGame,
			PrevGame1:   prevGame1,
			PrevGame2:   prevGame2,
			Bracket:     "L",
		}

		if prevGame1.Bracket == "L" {
			prevGame1.NextWinGame = game
		} else {
			prevGame1.NextLoseGame = game
		}
		if prevGame2.Bracket == "L" {
			prevGame2.NextWinGame = game
		} else {
			prevGame2.NextLoseGame = game
		}
	}

	if prevGame == 1 {
		nextWinGame.PrevGame1 = game
	}
	if prevGame == 2 {
		nextWinGame.PrevGame2 = game
	}
}

func (e *Elimination) setStartGames() {
	for i, game := range e.rounds[1] {
		if game.PrevGame1 != nil && game.PrevGame2 == nil {
			e.startGames[i] = game.PrevGame1
		} else if game.PrevGame1 == nil && game.PrevGame2 != nil {
			e.startGames[i] = game.PrevGame2
		} else if game.PrevGame1 != nil && game.PrevGame2 != nil {
			game1 := game.PrevGame1
			game2 := game.PrevGame2
			game3 := &Game{
				Round:     1,
				Bracket:   "L",
				PrevGame1: game1,
				PrevGame2: game2,
			}
			e.startGames[i] = game3
			game1.NextLoseGame = game3
			game2.NextLoseGame = game3
		} else {
			game.Bracket = "swap"
			e.startGames[i] = game
		}
	}
	jStart := len(e.startGames) - 1
	for i, game1 := range e.startGames {
		if game1.Bracket == "swap" {
			game1.Bracket = "W"

			for j := jStart; j > i; j-- {
				game2 := e.startGames[j]
				if game2.Bracket == "swap" {
					game2.Bracket = "W"

					e.startGames[i] = game2
					e.startGames[j] = game1
					jStart = j
					break
				}
			}
		}
	}
}

func (e *Elimination) numberGames(wGame *Game, lGame *Game, numberTwice bool) {
	wRound := 1
	lRound := 1

	e.order = 1

	e.numberGame(wGame, wRound, "W")
	wRound++

	e.numberGame(wGame, wRound, "W")
	wRound++

	e.numberGame(lGame, lRound, "L")
	lRound++

	if numberTwice {
		e.numberGame(lGame, lRound, "L")
		lRound++
	}

	for {
		checkFringe := e.numberGame(wGame, wRound, "W")
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

func (e *Elimination) numberGame(game *Game, round int, bracket string) string {
	if game == nil || game.Bracket != bracket || game.Order != 0 {
		return "null"
	}
	game1 := e.numberGame(game.PrevGame1, round, bracket)
	game2 := e.numberGame(game.PrevGame2, round, bracket)
	if game1 == "null" && game2 == "null" && game.Round == round {
		game.Order = e.order
		e.order++
		return "fringe"
	}
	return "not fringe"
}

func (e *Elimination) takeFirststartGames() *Game {
	game := e.startGames[e.startGamesPos]
	e.startGamesPos++
	return game
}

func (e *Elimination) takeFirstGame(round int) *Game {
	game := e.rounds[round-1][e.roundsPos[round-1]]
	e.roundsPos[round-1]++
	return game
}
