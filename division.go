package brackets

type Division struct {
	Teams []*Team
}

func (d *Division) String() string {
	s := "division: "
	s += "\n"
	s += "\tteams: ["

	s += "\n"
	for i := 0; i < len(d.Teams); i++ {
		s += "\t\t"
		s += d.Teams[i].String()
		s += "\n"
	}
	s += "\t]"

	s += "\n"
	return s
}
