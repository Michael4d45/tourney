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
	s := "\n" + tabs + "Game " + strconv.Itoa(g.GameNum) + "\n"
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

func DoubleGame(g double.Game, tabNum int) string {
	tabs := ""
	for i := 0; i < tabNum; i++ {
		tabs += "\t"
	}
	s := "\n" + tabs + "Game " + strconv.Itoa(g.GameNum) + "\n"
	s += tabs + "bracket: " + g.Bracket + "\n"
	s += tabs + "round: " + strconv.Itoa(g.Round) + "\n"
	if g.Team1 != nil {
		s += tabs + "team1: " + Team(*g.Team1) + "\n"
	}
	if g.Team2 != nil {
		s += tabs + "team2: " + Team(*g.Team2) + "\n"
	}

	if g.NextLoseGame != nil {
		s += tabs + "next lose: " + strconv.Itoa(g.NextLoseGame.GameNum) + "\n"
	}
	if g.NextWinGame != nil {
		s += tabs + "next win: " + strconv.Itoa(g.NextWinGame.GameNum) + "\n"
	}

	if g.PrevGame1 != nil {
		s += tabs + "prev1: " + DoubleGame(*g.PrevGame1, tabNum+1) + "\n"
	}
	if g.PrevGame2 != nil {
		s += tabs + "prev2: " + DoubleGame(*g.PrevGame2, tabNum+1) + "\n"
	}
	s += "\n"
	return s
}
