package brackets_test

import (
	"testing"

	"github.com/michael4d45/brackets"
)

func TestElimination(t *testing.T) {
	d := brackets.Division{
		Teams: []*brackets.Team{},
	}
	for i := 1; i <= 50; i++ {
		d.Teams = append(d.Teams, &brackets.Team{
			Seed: i,
		})
	}
	e := brackets.Elimination{}

	game := e.Generate(d)

    if game == nil {
        t.Error("Did not return game")
    }
}

func TestDoubleElimination(t *testing.T) {
	d := brackets.Division{
		Teams: []*brackets.Team{},
	}
	for i := 1; i <= 50; i++ {
		d.Teams = append(d.Teams, &brackets.Team{
			Seed: i,
		})
	}
	e := brackets.DoubleElimination{}

	game := e.Generate(d)

    if game == nil {
        t.Error("Did not return game")
    }
}

func BenchmarkElimination(b *testing.B) {
	d := brackets.Division{
		Teams: []*brackets.Team{},
	}
	for i := 1; i <= b.N; i++ {
		d.Teams = append(d.Teams, &brackets.Team{
			Seed: i,
		})
	}

    b.ResetTimer()
	e := brackets.Elimination{}

	e.Generate(d)
}

func BenchmarkDoubleElimination(b *testing.B) {
	d := brackets.Division{
		Teams: []*brackets.Team{},
	}
	for i := 1; i <= b.N; i++ {
		d.Teams = append(d.Teams, &brackets.Team{
			Seed: i,
		})
	}
    
    b.ResetTimer()
	e := brackets.DoubleElimination{}

	e.Generate(d)
}