package types

type HarrisPlayer struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Position  string `json:"pos"`
	Team      string `json:"team"`
	PPRRank   int    `json:"pprRank"`
	StdRank   int    `json:"stdRank"`
}

type HarrisPlayerMatch struct {
	Harris *HarrisPlayer `json:"harrisPlayer"`
	Player *Player       `json:"player"`
}

func (m *HarrisPlayerMatch) AddPlayerRank() {
	if m.Player != nil && m.Harris != nil {
		m.Player.CustomStdRank = m.Harris.StdRank
		m.Player.CustomPprRank = m.Harris.PPRRank
	}

	return
}
