package brackets

import (
	"math"
	"sort"
)

type DoubleElimination struct {
	gg        [][]*Game
	ggPos     []int
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

func (d *DoubleElimination) Generate(division Division) *Game {
	e := Elimination{}
	wGame := e.Generate(division)
	if wGame == nil {
		return nil
	}

	count := len(division.Teams)
	lRound := loserBracketRoundCount(count)
	wRound := wGame.round
	d.gg = make([][]*Game, wRound)
	d.ggPos = make([]int, wRound)
	d.addGames(wGame)
	for i := 0; i < wRound; i++ {
		lowToHigh := (i % 2) != 0
		sort.Slice(d.gg[i], func(i, j int) bool {
			if lowToHigh {
				return i < j
			} else {
				return j < i
			}
		})
	}

	ifFirstLost := &Game{
		bracket: "WW",
		round:   lRound,
	}
	lRound--

	wlGame := &Game{
		bracket:      "WL",
		round:        lRound,
		nextWinGame:  ifFirstLost,
		nextLoseGame: ifFirstLost,
		prevGame1:    wGame,
	}
	lRound--

	ifFirstLost.prevGame1 = wlGame
	ifFirstLost.prevGame2 = wlGame

	if wRound > 1 {
		d.bucket = []*Game{}
		d.setBucket()

		d.loserBracket(wlGame, 2, lRound, wRound)
	}

	lGame := wlGame.prevGame2

	n := math.Log2(float64(count))
	o := int(math.Pow(2, n))
	p := o - (o / 4)
	numberTwice := (count <= o) && (count > p)

	d.numberGames(ifFirstLost, wlGame, wGame, lGame, numberTwice)

	return ifFirstLost
}

func (d *DoubleElimination) loserBracket(nextWinGame *Game, prevGame int, lRound int, wRound int) {
	if wRound == 2 {
		d.firstRoundGames(nextWinGame, prevGame, lRound)
		return
	}

	wGame := d.takeFirstGame(wRound)
	fromWinnerGame := &Game{
		round:       lRound,
		nextWinGame: nextWinGame,
		prevGame1:   wGame,
		bracket:     "L",
	}
	lRound--

	wGame.nextLoseGame = fromWinnerGame
	if prevGame == 1 {
		nextWinGame.prevGame1 = fromWinnerGame
	}
	if prevGame == 2 {
		nextWinGame.prevGame2 = fromWinnerGame
	}

	losersPlayGame := &Game{
		round:       lRound,
		nextWinGame: fromWinnerGame,
		bracket:     "L",
	}
	lRound--

	fromWinnerGame.prevGame2 = losersPlayGame

	wRound--
	d.loserBracket(losersPlayGame, 1, lRound, wRound)
	d.loserBracket(losersPlayGame, 2, lRound, wRound)
}

func (d *DoubleElimination) firstRoundGames(nextWinGame *Game, prevGame int, lRound int) {
	prevGame1 := d.takeFirstBucket()
	prevGame2 := d.takeFirstGame(2)
	game := prevGame2
	if prevGame1.round != 2 {
		game := &Game{
			round:       lRound,
			nextWinGame: nextWinGame,
			prevGame1:   prevGame1,
			prevGame2:   prevGame2,
			bracket:     "L",
		}
		lRound--
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

func (d *DoubleElimination) setBucket() {
	round := 2
	for i := 0; i < len(d.gg[round-1]); i++ {
		game := d.gg[round-1][i]
		if game.prevGame1 != nil && game.prevGame2 == nil {
			d.bucket = append(d.bucket, game.prevGame1)
		} else if game.prevGame1 == nil && game.prevGame2 != nil {
			d.bucket = append(d.bucket, game.prevGame2)
		} else if game.prevGame1 != nil && game.prevGame2 != nil {
			game1 := game.prevGame1
			game2 := game.prevGame2
			game3 := &Game{
				round:     1,
				bracket:   "L",
				prevGame1: game1,
				prevGame2: game2,
			}
			d.bucket = append(d.bucket, game3)
			game1.nextLoseGame = game3
			game2.nextLoseGame = game3
		} else {
			game.bracket = "swap"
			d.bucket = append(d.bucket, game)
		}
	}
	jStart := len(d.bucket) - 1
	for i, game1 := range d.bucket {
		if game1.bracket == "swap" {
			for j := jStart; j > i; j-- {
				game2 := d.bucket[j]
				if game2.bracket == "swap" {
					d.bucket[i] = game2
					d.bucket[j] = game1
					game2.bracket = "W"
					game1.bracket = "W"
					jStart = j
					break
				}
			}
			game1.bracket = "W"
		}
	}
}

func (d *DoubleElimination) numberGames(lastGame *Game, wlGame *Game, wGame *Game, lGame *Game, numberTwice bool) {
	wRound := 1
	lRound := 1

	d.gameNum = 1
	d.numberGame(wGame, wRound, "W")
	wRound++
	d.numberGame(wGame, wRound, "W")
	wRound++
	d.numberGame(lGame, lRound, "L")
	lRound++
	if numberTwice {
		d.numberGame(lGame, lRound, "L")
		lRound++
	}

	for ; d.numberGame(wGame, wRound, "W") == "not fringe"; wRound++ {

		d.numberGame(lGame, lRound, "L")
		lRound++

		d.numberGame(lGame, lRound, "L")
		lRound++
	}
	for ; d.numberGame(lGame, lRound, "L") == "not fringe"; lRound++ {

	}
	wlGame.gameNum = d.gameNum
	d.gameNum++
	lastGame.gameNum = d.gameNum
	d.gameNum++
}

func (d *DoubleElimination) numberGame(game *Game, round int, bracket string) string {
	if game == nil || game.bracket != bracket || game.gameNum != 0 {
		return "null"
	}
	game1 := d.numberGame(game.prevGame1, round, bracket)
	game2 := d.numberGame(game.prevGame2, round, bracket)
	if game1 == "null" && game2 == "null" && game.round == round {
		game.gameNum = d.gameNum
		d.gameNum++
		return "fringe"
	}
	return "not fringe"
}

func (d *DoubleElimination) addGames(game *Game) {
	if game == nil {
		return
	}
	d.addGame(game)
	d.addGames(game.prevGame1)
	d.addGames(game.prevGame2)
}

func (d *DoubleElimination) addGame(game *Game) {
	if game == nil {
		return
	}
	game.bracket = "W"
	d.gg[game.round-1] = append(d.gg[game.round-1], game)
}

func (d *DoubleElimination) takeFirstBucket() *Game {
	game := d.bucket[d.bucketPos]
	d.bucketPos++
	return game
}

func (d *DoubleElimination) takeFirstGame(round int) *Game {
	game := d.gg[round-1][d.ggPos[round-1]]
	d.ggPos[round-1]++
	return game
}
