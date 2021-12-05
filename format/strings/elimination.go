package strings

import (
	"strconv"

	"github.com/michael4d45/tourney/elim"
)

func Elim(g elim.Games) string {
	return elimGame(*g.FinalGame, 0, make(map[elim.Game]struct{}))
}

func shortGame(g elim.Game) string {
	return "Game " + strconv.Itoa(g.Order) + " : " + g.Bracket + "\n"
}

func elimGame(g elim.Game, tabNum int, games map[elim.Game]struct{}) string {
	_, exists := games[g]
	if exists {
		return shortGame(g)
	}
	games[g] = struct{}{}

	tabs := ""
	for i := 0; i < tabNum; i++ {
		tabs += "\t"
	}

	s := "\n" + tabs + "Game " + strconv.Itoa(g.Order) + " : " + g.Bracket + "\n"
	s += tabs + "round: " + strconv.Itoa(g.Round) + "\n"
	if g.NextLoseGame != nil {
		s += tabs + "next lose: " + shortGame(*g.NextLoseGame)
	}
	if g.NextWinGame != nil {
		s += tabs + "next win: " + shortGame(*g.NextWinGame)
	}

	if g.Team1 != nil {
		s += tabs + "team1: " + Team(*g.Team1) + "\n"
	}
	if g.PrevGame1 != nil {
		s += tabs + "prev1: " + elimGame(*g.PrevGame1, tabNum+1, games) + "\n"
	}

	if g.Team2 != nil {
		s += tabs + "team2: " + Team(*g.Team2) + "\n"
	}
	if g.PrevGame2 != nil {
		s += tabs + "prev2: " + elimGame(*g.PrevGame2, tabNum+1, games) + "\n"
	}
	s += "\n"
	return s
}
