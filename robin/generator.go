package robin

import "github.com/michael4d45/tourney"

type Games struct {
	Rounds [][]*Game
}

func Generate(d tourney.Division) Games {
	isOdd := len(d.Teams)%2 != 0

	if isOdd {
		d.Teams = append([]*tourney.Team{nil}, d.Teams...)
	}

	numRounds := len(d.Teams) - 1
	numGames := len(d.Teams) / 2

	games := Games{
		Rounds: make([][]*Game, numRounds),
	}

	for i := range games.Rounds {
		games.Rounds[i] = make([]*Game, numGames)
	}

	games.circle(d)

	if isOdd {
		for i, rr := range games.Rounds {
			games.Rounds[i] = rr[1:]
		}
	}

	return games
}
