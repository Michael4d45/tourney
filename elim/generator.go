package elim

import (
	"sort"

	"github.com/michael4d45/tourney"
)

type Games struct {
	Kind string

	FinalGame *Game
}

// [even/odd],[lower/higher], random
func Generate(division tourney.Division, topOrder string, elimType string) Games {
	games := Games{
		Kind: elimType,
	}

	if elimType == "single" || elimType == "double" {
		s := &singleGenerator{}
		g := s.generateSingle(division, topOrder)

		if elimType == "double" && g != nil {
			l := &loserGen{}
			g = l.generateLoser(g, len(division.Teams), s.rounds)
		}

		games.FinalGame = g
	}

	return games
}

func Rounds(g *Game) [][]*Game {
	rounds := make([][]*Game, g.Round)
	gamesMap := genGamesMap(g, map[*Game]struct{}{})

	for gg := range gamesMap {
		rounds[gg.Round-1] = append(rounds[gg.Round-1], gg)
	}

	for _, gg := range rounds {
		sort.Slice(gg, func(i, j int) bool {
			return gg[i].Order < gg[j].Order
		})
	}

	return rounds
}

func genGamesMap(g *Game, games map[*Game]struct{}) map[*Game]struct{} {
	if _, exists := games[g]; exists || g == nil {
		return games
	}
	games[g] = struct{}{}
	games = genGamesMap(g.PrevGame1, games)
	games = genGamesMap(g.PrevGame2, games)
	return games
}