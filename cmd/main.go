package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elim"

	"github.com/michael4d45/tourney/format/strings"
)

func main() {
	var printTime bool
	var bracket string
	var printBracket bool
	var teamCount int
	var topOrder string

	flag.BoolVar(&printTime, "t", false, "output time")
	flag.BoolVar(&printBracket, "p", false, "output bracket")
	flag.StringVar(&bracket, "b", "double", "Type of bracket; double or single")
	flag.IntVar(&teamCount, "c", 10, "Number of teams")
	flag.StringVar(&topOrder, "o", "odd", "How the games should print out")

	flag.Parse()

	d := tourney.Division{}
	d.MakeTeams(teamCount)

	if printBracket {
		fmt.Println(strings.Division(d))
	}
	start := time.Now()

	games := elim.Generate(d, topOrder, bracket)
	if printBracket {
		fmt.Println(strings.Elim(games))
	}

	if printTime {
		elapsed := time.Since(start)
		fmt.Printf("took %s\n", elapsed)
	}
}
