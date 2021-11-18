package main

import (
	"fmt"
	"time"

	"github.com/michael4d45/brackets"
)

func main() {
	d := brackets.Division{
		Teams: []*brackets.Team{},
	}
	for i := 1; i <= 10; i++ {
		d.Teams = append(d.Teams, &brackets.Team{
			Seed: i,
			// Division: &d,
		})
	}
	fmt.Println(d.String())
	
	start := time.Now()

	elim := brackets.Elimination{}

	elim.Generate(d)
	// e_game := elim.Generate(d)

	// fmt.Println(e_game.String(0))

	e := brackets.DoubleElimination{}

	e.Generate(d)
	// game := e.Generate(d)

	// fmt.Println(game.String(0))

	elapsed := time.Since(start)

	fmt.Printf("took %s\n", elapsed)
}
