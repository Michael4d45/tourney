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
	bucket    []*Game
	bucketPos int
	gameNum   int
}

func loserBracketRoundCount(n int) int {
	f := float64(n - 1)
	o := math.Log2(f)
	p := math.Log2(f * (8.0 / 3.0))
	return int(o) + int(p)
}

func shouldNumberTwice(count int) bool {
	n := math.Floor(math.Log2(float64(count)))
	o := math.Pow(2, n)
	p := o - (o / 4)
	c := float64(count)
	return (c <= o) && (c > p)
}

func (e *Elimination) Generate(division tourney.Division) *Game {
	elim := single.Elimination{}
	eGame := elim.Generate(division)
	if eGame == nil {
		return nil
	}

	count := len(division.Teams)
	lRound := loserBracketRoundCount(count)

	ifFirstLost := &Game{
		Bracket: "WW",
		Round:   lRound,
	}
	lRound--

	wlGame := &Game{
		Bracket:      "WL",
		Round:        lRound,
		NextWinGame:  ifFirstLost,
		NextLoseGame: ifFirstLost,
	}
	lRound--

	wGame, gamesMap := copy(eGame, map[*single.Game]*Game{})

	wGame.NextWinGame = wlGame

	ifFirstLost.PrevGame1 = wlGame
	ifFirstLost.PrevGame2 = wlGame

	wlGame.PrevGame1 = wGame

	wRound := wGame.Round
	e.rounds = make([][]*Game, wRound)
	e.roundsPos = make([]int, wRound)

	for _, gg := range gamesMap {
		e.rounds[gg.Round-1] = append(e.rounds[gg.Round-1], gg)
	}

	if wRound > 1 {
		sort.Slice(e.rounds[1], func(j, k int) bool {
			return e.rounds[1][j].GameNum < e.rounds[1][k].GameNum
		})
		e.bucket = make([]*Game, len(e.rounds[1]))
		e.setBucket()

		for i, rr := range e.rounds {
			lowToHigh := (i % 2) != 0
			sort.Slice(rr, func(j, k int) bool {
				if lowToHigh {
					return rr[k].GameNum < rr[j].GameNum
				} else {
					return rr[j].GameNum < rr[k].GameNum
				}
			})
		}

		e.loserBracket(wlGame, 2, lRound, wRound)
	}

	lGame := wlGame.PrevGame2

	numberTwice := shouldNumberTwice(count)

	for _, gg := range gamesMap {
		gg.Bracket = "W"
		gg.GameNum = 0
	}

	e.numberGames(wGame, lGame, numberTwice)
	wlGame.GameNum = e.gameNum
	e.gameNum++
	ifFirstLost.GameNum = e.gameNum
	e.gameNum++

	return ifFirstLost
}

func (e *Elimination) loserBracket(nextWinGame *Game, prevGame int, lRound int, wRound int) {
	if wRound == 2 {
		e.firstRoundGames(nextWinGame, prevGame, lRound)
		return
	}

	wGame := e.takeFirstGame(wRound)
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
	prevGame1 := e.takeFirstBucket()
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

func (e *Elimination) setBucket() {
	for i, game := range e.rounds[1] {
		if game.PrevGame1 != nil && game.PrevGame2 == nil {
			e.bucket[i] = game.PrevGame1
		} else if game.PrevGame1 == nil && game.PrevGame2 != nil {
			e.bucket[i] = game.PrevGame2
		} else if game.PrevGame1 != nil && game.PrevGame2 != nil {
			game1 := game.PrevGame1
			game2 := game.PrevGame2
			game3 := &Game{
				Round:     1,
				Bracket:   "L",
				PrevGame1: game1,
				PrevGame2: game2,
			}
			e.bucket[i] = game3
			game1.NextLoseGame = game3
			game2.NextLoseGame = game3
		} else {
			game.Bracket = "swap"
			e.bucket[i] = game
		}
	}
	jStart := len(e.bucket) - 1
	for i, game1 := range e.bucket {
		if game1.Bracket == "swap" {
			game1.Bracket = "W"

			for j := jStart; j > i; j-- {
				game2 := e.bucket[j]
				if game2.Bracket == "swap" {
					game2.Bracket = "W"

					e.bucket[i] = game2
					e.bucket[j] = game1
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

	e.gameNum = 1

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
	if game == nil || game.Bracket != bracket || game.GameNum != 0 {
		return "null"
	}
	game1 := e.numberGame(game.PrevGame1, round, bracket)
	game2 := e.numberGame(game.PrevGame2, round, bracket)
	if game1 == "null" && game2 == "null" && game.Round == round {
		game.GameNum = e.gameNum
		e.gameNum++
		return "fringe"
	}
	return "not fringe"
}

func (e *Elimination) takeFirstBucket() *Game {
	game := e.bucket[e.bucketPos]
	e.bucketPos++
	return game
}

func (e *Elimination) takeFirstGame(round int) *Game {
	game := e.rounds[round-1][e.roundsPos[round-1]]
	e.roundsPos[round-1]++
	return game
}
