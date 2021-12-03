package elim

import (
	"math"
	"sort"
)

type loserGen struct {
	rounds    [][]*Game
	roundsPos []int

	startGames    []*Game
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
func (l *loserGen) generateLoser(game *Game, teamCount int, rounds [][]*Game) *Game {
	l.rounds = rounds

	lRound := loserBracketRoundCount(teamCount)

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

	game.NextWinGame = wlGame

	wlGame.PrevGame1 = game

	wRound := game.Round
	l.roundsPos = make([]int, wRound)

	if wRound > 1 {

		// Collect starting games to hook together the winner's
		// bracket to the looser's bracket.
		l.setStartGames()

		// Sort the games in order to get a deterministic
		// winner to loser connection.
		for i, rr := range l.rounds {
			lowToHigh := (i % 2) != 0
			sort.Slice(rr, func(j, k int) bool {
				if lowToHigh {
					return rr[k].Order < rr[j].Order
				} else {
					return rr[j].Order < rr[k].Order
				}
			})
		}

		// Generate the loosers bracket.
		l.loserBracket(wlGame, 2, lRound, wRound)
	}

	// Set numbers on the games.
	lGame := wlGame.PrevGame2

	for _, gg := range l.rounds {
		for _, g := range gg {
			g.Bracket = "W"
			g.Order = 0
		}
	}

	numberTwice := shouldNumberTwice(teamCount)
	l.numberGames(game, lGame, numberTwice)
	wlGame.Order = l.order
	l.order++
	ifFirstLost.Order = l.order
	l.order++

	return ifFirstLost
}

// loserBracket generates the losers bracket.
func (l *loserGen) loserBracket(nextWinGame *Game, prevGame int, lRound int, wRound int) {
	// Use different rules when on the first round.
	if wRound == 2 {
		l.firstRoundGames(nextWinGame, prevGame, lRound)
		return
	}

	// Get a game from the winners bracket.
	game := l.takeFirstGame(wRound)

	// Create a loser bracket game that connects to
	// the winners bracket.
	fromWinnerGame := &Game{
		Round:       lRound,
		NextWinGame: nextWinGame,
		PrevGame1:   game,
		Bracket:     "L",
	}
	lRound--

	game.NextLoseGame = fromWinnerGame

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

	l.loserBracket(losersPlayGame, 1, lRound, wRound-1)
	l.loserBracket(losersPlayGame, 2, lRound, wRound-1)
}

func (l *loserGen) firstRoundGames(nextGame *Game, prevGame int, lRound int) {
	prevGame1 := l.takeFirstStartGames()
	prevGame2 := l.takeFirstGame(2)
	var game *Game

	if prevGame1.Round != 2 {
		game = &Game{
			Round:       lRound,
			NextWinGame: nextGame,
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
	} else {
		game = prevGame2

		if nextGame.Bracket == "L" {
			game.NextLoseGame = nextGame
		} else {
			game.NextWinGame = nextGame
		}
	}

	if prevGame == 1 {
		nextGame.PrevGame1 = game
	}
	if prevGame == 2 {
		nextGame.PrevGame2 = game
	}
}

func (l *loserGen) setStartGames() {
	sort.Slice(l.rounds[1], func(j, k int) bool {
		return l.rounds[1][j].Order < l.rounds[1][k].Order
	})

	l.startGames = make([]*Game, len(l.rounds[1]))

	for i, game := range l.rounds[1] {
		if game.PrevGame1 != nil && game.PrevGame2 == nil {
			l.startGames[i] = game.PrevGame1
		} else if game.PrevGame1 == nil && game.PrevGame2 != nil {
			l.startGames[i] = game.PrevGame2
		} else if game.PrevGame1 != nil && game.PrevGame2 != nil {
			game1 := game.PrevGame1
			game2 := game.PrevGame2
			game3 := &Game{
				Round:     1,
				PrevGame1: game1,
				PrevGame2: game2,
				Bracket:   "L",
			}
			l.startGames[i] = game3
			game1.NextLoseGame = game3
			game2.NextLoseGame = game3
		} else {
			game.Bracket = "swap"
			l.startGames[i] = game
		}
	}

	jStart := len(l.startGames) - 1
	for i, game1 := range l.startGames {
		if game1.Bracket == "swap" {
			game1.Bracket = "W"

			for j := jStart; j > i; j-- {
				game2 := l.startGames[j]
				if game2.Bracket == "swap" {
					game2.Bracket = "W"

					l.startGames[i] = game2
					l.startGames[j] = game1
					jStart = j
					break
				}
			}
		}
	}
}

func (l *loserGen) numberGames(game *Game, lGame *Game, numberTwice bool) {
	wRound := 1
	lRound := 1

	l.order = 1

	l.numberGame(game, wRound, "W")
	wRound++

	l.numberGame(game, wRound, "W")
	wRound++

	l.numberGame(lGame, lRound, "L")
	lRound++

	if numberTwice {
		l.numberGame(lGame, lRound, "L")
		lRound++
	}

	for {
		checkFringe := l.numberGame(game, wRound, "W")
		wRound++

		if checkFringe != "not fringe" {
			break
		}

		l.numberGame(lGame, lRound, "L")
		lRound++

		l.numberGame(lGame, lRound, "L")
		lRound++
	}

	for {
		checkFringe := l.numberGame(lGame, lRound, "L")
		lRound++
		if checkFringe != "not fringe" {
			break
		}
	}
}

func (l *loserGen) numberGame(game *Game, round int, bracket string) string {
	if game == nil || game.Bracket != bracket || game.Order != 0 {
		return "null"
	}
	game1 := l.numberGame(game.PrevGame1, round, bracket)
	game2 := l.numberGame(game.PrevGame2, round, bracket)
	if game1 == "null" && game2 == "null" && game.Round == round {
		game.Order = l.order
		l.order++
		return "fringe"
	}
	return "not fringe"
}

func (l *loserGen) takeFirstStartGames() *Game {
	game := l.startGames[l.startGamesPos]
	l.startGamesPos++
	return game
}

func (l *loserGen) takeFirstGame(round int) *Game {
	game := l.rounds[round-1][l.roundsPos[round-1]]
	l.roundsPos[round-1]++
	return game
}
