package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/michael4d45/tourney"
	"github.com/michael4d45/tourney/elim/double"
	"github.com/michael4d45/tourney/elim/single"

	"github.com/michael4d45/tourney/elim"

	"github.com/michael4d45/tourney/format/strings"
)

func main() {
	var printTime bool
	var bracket string
	var printBracket bool
	var teamCount int

	flag.BoolVar(&printTime, "t", false, "output time")
	flag.BoolVar(&printBracket, "p", false, "output bracket")
	flag.StringVar(&bracket, "b", "double", "Type of bracket; double or single")
	flag.IntVar(&teamCount, "c", 10, "Number of teams")

	flag.Parse()

	d := tourney.Division{}
	d.MakeTeams(teamCount)
	
	if printBracket {
		fmt.Println(strings.Division(d))
	}
	start := time.Now()

	switch bracket {
	case "single":
		gen := single.Elimination{}
		game := gen.Generate(d)
		if printBracket {
			fmt.Println(strings.SingleGame(*game, 0))
		}
	case "double":
		gen := double.Elimination{}
		game := gen.Generate(d)
		if printBracket {
			fmt.Println(strings.DoubleGame(*game, 0, map[double.Game]struct{}{}))
		}
	case "elim":
		game := elim.GenerateDouble(d)
		if printBracket {
			fmt.Println(strings.DoublesGame(game, 0, map[elim.Game]struct{}{}))
		}
	}

	if printTime {
		elapsed := time.Since(start)
		fmt.Printf("took %s\n", elapsed)
	}
}
