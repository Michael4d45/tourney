package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/michael4d45/brackets"
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

	d := brackets.Division{
		Teams: []*brackets.Team{},
	}
	for i := 1; i <= teamCount; i++ {
		d.Teams = append(d.Teams, &brackets.Team{
			Seed: i,
			// Division: &d,
		})
	}
	fmt.Println(d.String())

	start := time.Now()

	var game *brackets.Game

	switch bracket {
	case "single":
		e := brackets.Elimination{}
		game = e.Generate(d)
	case "double":
		e := brackets.DoubleElimination{}
		game = e.Generate(d)
	}

	if printBracket {
		fmt.Println(game.String(0))
	}

	if printTime {
		elapsed := time.Since(start)
		fmt.Printf("took %s\n", elapsed)
	}
}
