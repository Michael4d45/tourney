package robin

import "github.com/michael4d45/tourney"

type Games struct {
	Rounds [][]*Game
}

func Generate(d tourney.Division) Games {
	var numRounds int
	var numGames int
	if len(d.Teams)%2 == 0 {
		numRounds = len(d.Teams) - 1
		numGames = len(d.Teams) / 2
	} else {
		numRounds = len(d.Teams)
		numGames = (len(d.Teams) - 1) / 2
	}
	games := Games{
		Rounds: make([][]*Game, numRounds),
	}

	for i := range games.Rounds {
		games.Rounds[i] = make([]*Game, numGames)
	}

	games.circle(d)

	return games
}
