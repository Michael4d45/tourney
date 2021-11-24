package strings

import (
	"strconv"

	"github.com/michael4d45/tourney/elimination/double"
	"github.com/michael4d45/tourney/elimination/single"
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
