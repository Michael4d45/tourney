package strings

import (
	"strconv"

	"github.com/michael4d45/tourney/elim"
	"github.com/michael4d45/tourney/elim/double"
	"github.com/michael4d45/tourney/elim/single"
)

func SingleGame(g single.Game, tabNum int) string {
	tabs := ""
	for i := 0; i < tabNum; i++ {
		tabs += "\t"
	}
	s := "\n" + tabs + "Game " + strconv.Itoa(g.Order) + "\n"
	s += tabs + "round: " + strconv.Itoa(g.Round) + "\n"
	if g.Team1 != nil {
		s += tabs + "team1: " + Team(*g.Team1) + "\n"
	}
	if g.Team2 != nil {
		s += tabs + "team2: " + Team(*g.Team2) + "\n"
	}

	if g.PrevGame1 != nil {
		s += tabs + "prev1: " + SingleGame(*g.PrevGame1, tabNum+1) + "\n"
	}
	if g.PrevGame2 != nil {
		s += tabs + "prev2: " + SingleGame(*g.PrevGame2, tabNum+1) + "\n"
	}
	s += "\n"
	return s
}

func shortDoubleGame(g double.Game) string {
	return "Game " + strconv.Itoa(g.Order) + " : " + g.Bracket + "\n"
}

func DoubleGame(g double.Game, tabNum int, games map[double.Game]struct{}) string {
	_, exists := games[g]
	if exists {
		return shortDoubleGame(g)
	}
	games[g] = struct{}{}

	tabs := ""
	for i := 0; i < tabNum; i++ {
		tabs += "\t"
	}

	s := "\n" + tabs + "Game " + strconv.Itoa(g.Order) + " : " + g.Bracket + "\n"
	s += tabs + "round: " + strconv.Itoa(g.Round) + "\n"
	if g.Team1 != nil {
		s += tabs + "team1: " + Team(*g.Team1) + "\n"
	}
	if g.Team2 != nil {
		s += tabs + "team2: " + Team(*g.Team2) + "\n"
	}

	if g.NextLoseGame != nil {
		s += tabs + "next lose: " + shortDoubleGame(*g.NextLoseGame)
	}
	if g.NextWinGame != nil {
		s += tabs + "next win: " + shortDoubleGame(*g.NextWinGame)
	}

	if g.PrevGame1 != nil {
		s += tabs + "prev1: " + DoubleGame(*g.PrevGame1, tabNum+1, games) + "\n"
	}
	if g.PrevGame2 != nil {
		s += tabs + "prev2: " + DoubleGame(*g.PrevGame2, tabNum+1, games) + "\n"
	}
	s += "\n"
	return s
}

func shortDoublesGame(g elim.Game) string {
	return "Game " + strconv.Itoa(g.Order()) + "\n" //+ " : " + g.Bracket() + "\n"
}

func DoublesGame(g elim.Game, tabNum int, games map[elim.Game]struct{}) string {
	if elim.IsElimNil(g) {
		return "nil"
	}
	_, exists := games[g]
	if exists {
		return shortDoublesGame(g)
	}
	games[g] = struct{}{}

	tabs := ""
	for i := 0; i < tabNum; i++ {
		tabs += "\t"
	}

	team1, team2 := g.Teams()
	s := "\n" + tabs + shortDoublesGame(g)
	s += tabs + "round: " + strconv.Itoa(g.Round()) + "\n"
	if team1 != nil {
		s += tabs + "team1: " + Team(*team1) + "\n"
	}
	if team2 != nil {
		s += tabs + "team2: " + Team(*team2) + "\n"
	}

	// nextLoseGame := g.NextLoseGame()
	// nextWinGame := g.NextWinGame()

	// if nextLoseGame != nil {
	// 	s += tabs + "next lose: " + shortDoublesGame(*nextLoseGame)
	// }
	// if nextWinGame != nil {
	// 	s += tabs + "next win: " + shortDoublesGame(*nextWinGame)
	// }

	if g.PrevGame1() != nil {
		s += tabs + "prev1: " + DoublesGame(*g.PrevGame1(), tabNum+1, games) + "\n"
	}
	if g.PrevGame2() != nil {
		s += tabs + "prev2: " + DoublesGame(*g.PrevGame2(), tabNum+1, games) + "\n"
	}
	s += "\n"
	return s
}
