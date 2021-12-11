package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elim"
	"github.com/michael4d45/tourney/robin"

	"github.com/michael4d45/tourney/format/strings"
)

func main() {
	var printTime bool
	var bracket string
	var printBracket bool
	var teamCount int
	var topOrder string
	var printAnalytics bool

	flag.BoolVar(&printTime, "t", false, "output time")
	flag.BoolVar(&printBracket, "p", false, "output bracket")
	flag.StringVar(&bracket, "b", "double", "Type of bracket; double or single")
	flag.IntVar(&teamCount, "c", 10, "Number of teams")
	flag.StringVar(&topOrder, "o", "odd", "How the games should print out")
	flag.BoolVar(&printAnalytics, "a", false, "output numbers")

	flag.Parse()

	d := tourney.Division{}
	d.MakeTeams(teamCount)

	if printBracket {
		fmt.Println(strings.Division(d))
	}
	start := time.Now()

	if bracket == "double" || bracket == "single" {
		games := elim.Generate(d, topOrder, bracket)
		if printBracket {
			fmt.Println(strings.Elim(games))
		}
		if printAnalytics {
			rounds := elim.Rounds(games.FinalGame)
			fmt.Println(len(rounds), "rounds")
			var numGames int
			for _, rr := range rounds {
				numGames += len(rr)
			}
			fmt.Println(numGames, "games")
		}
	}
	if bracket == "robin" {
		games := robin.Generate(d)
		if printBracket {
			fmt.Println(strings.Robin(games))
		}
		if printAnalytics {
			fmt.Println(len(games.Rounds), "rounds")
			var numGames int
			for _, rr := range games.Rounds {
				numGames += len(rr)
			}
			fmt.Println(numGames, "games")
		}
	}

	if printTime {
		elapsed := time.Since(start)
		fmt.Printf("took %s\n", elapsed)
	}
}
