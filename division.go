package tourney

type Division struct {
	Teams []*Team
}

func (d *Division) MakeTeams(num int) {
	d.Teams = make([]*Team, num)
	for i := range d.Teams {
		d.Teams[i] = &Team{
			Seed: i + 1,
		}
	}
}
