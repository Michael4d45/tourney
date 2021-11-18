package brackets

import (
	"strconv"
)

type Team struct {
	// Division *Division
	Seed int
	Name string
}

func (t *Team) String() string {
	if t == nil {
		return "nil"
	}
	s := "Team: "
	s += "seed: " + strconv.Itoa(t.Seed)
	// s += ", d: " + strconv.Itoa(len(t.Division.Teams))
	return s
}
