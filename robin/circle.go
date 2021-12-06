package robin

import "github.com/michael4d45/tourney"

func (g *Games) circle(d tourney.Division) {
	numTeams := len(d.Teams)
	for i, rr := range g.Rounds {
		rr[0] = &Game{
			Team1: d.Teams[0],
			Team2: d.Teams[team2(numTeams, i, 0)],
		}
		for j := 1; j < len(rr); j++ {
			rr[j] = &Game{
				Team1: d.Teams[team1(numTeams, i, j)],
				Team2: d.Teams[team2(numTeams, i, j)],
			}
		}
	}
}

func team1(numTeams, i, j int) int {
	team1 := j - i
	return normTeamNum(numTeams, team1)
}

func team2(numTeams, i, j int) int {
	team2 := (numTeams - 1) - j - i
	return normTeamNum(numTeams, team2)
}

func normTeamNum(numTeams, n int) int {
	n %= numTeams
	if n <= 0 {
		n += numTeams - 1
	}
	return n
}
