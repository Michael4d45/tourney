package brackets

import "strconv"

type Game struct {
	team1   *Team
	team2   *Team
	round   int
	gameNum int

	nextWinGame  *Game
	nextLoseGame *Game
	prevGame1    *Game
	prevGame2    *Game

	bracket string
}

func (g *Game) String(tabNum int) string {
	tabs := ""
	if g == nil {
		return "nil"
	}
	for i := 0; i < tabNum; i++ {
		tabs += "\t"
	}
	s := "\n" + tabs + "Game " + strconv.Itoa(g.gameNum) + "\n"
	s += tabs + "bracket: " + g.bracket + "\n"
	s += tabs + "round: " + strconv.Itoa(g.round) + "\n"
	s += tabs + "team1: " + g.team1.String() + "\n"
	s += tabs + "team2: " + g.team2.String() + "\n"

	if g.nextLoseGame != nil {
		s += tabs + "next lose: " + strconv.Itoa(g.nextLoseGame.gameNum) + "\n"

	}
	if g.nextWinGame != nil {
		s += tabs + "next win: " + strconv.Itoa(g.nextWinGame.gameNum) + "\n"
	}

	s += tabs + "prev1: " + g.prevGame1.String(tabNum+1) + "\n"
	s += tabs + "prev2: " + g.prevGame2.String(tabNum+1) + "\n"
	s += "\n"
	return s
}
