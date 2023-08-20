package types

import (
	"errors"
	"strconv"
)

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

type HarrisEspnPlayerMatch struct {
	Harris *HarrisPlayer `json:"harrisPlayer"`
	Espn   *EspnPlayer   `json:"espnPlayer"`
}

func (m HarrisEspnPlayerMatch) ToPlayer() (out *Player, err error) {
	if m.Harris == nil && m.Espn == nil {
		return nil, errors.New("Both Harris and ESPN ranks are not found")
	}
	if m.Espn != nil {
		out = &Player{
			Id:                strconv.Itoa(m.Espn.Id),
			FirstName:         m.Espn.Profile.FirstName,
			LastName:          m.Espn.Profile.LastName,
			Name:              m.Espn.Profile.FullName,
			Position:          m.Espn.Position(),
			Team:              string(m.Espn.Position()),
			EspnOvrStdRank:    m.Espn.Profile.Ranks.Standard.Rank,
			EspnOvrPprRank:    m.Espn.Profile.Ranks.PPR.Rank,
			EspnPlayerOutlook: m.Espn.Profile.SeasonOutlook,
			SeasonStats:       []*Stats{},
		}
		if m.Espn.Profile.Ownership != nil {
			out.EspnAdp = m.Espn.Profile.Ownership.AvgDraftPos
		}
		if m.Espn.Profile.StatsBySeason != nil {
			for _, espnSsnStats := range m.Espn.Profile.StatsBySeason {
				ssnStats := &Stats{
					PPG:    espnSsnStats.AvgPPGame,
					MinPPG: espnSsnStats.AvgPPGame,
					MaxPPG: espnSsnStats.AvgPPGame,
					Year:   espnSsnStats.Season,
				}
				if espnSsnStats.Stats != nil {
					ssnStats.GamesPlayed = espnSsnStats.Stats.GamesPlayed
					ssnStats.RushAttempts = espnSsnStats.Stats.RushAttempts
					ssnStats.RushYards = espnSsnStats.Stats.RushYards
					ssnStats.RushTds = espnSsnStats.Stats.RushTds
					ssnStats.Recs = espnSsnStats.Stats.Recs
					ssnStats.RecYards = espnSsnStats.Stats.RecYards
					ssnStats.RecTds = espnSsnStats.Stats.RecTds
					ssnStats.PassAttempts = espnSsnStats.Stats.PassAttempts
					ssnStats.PassCompletions = espnSsnStats.Stats.PassCompletions
					ssnStats.PassYards = espnSsnStats.Stats.PassYards
					ssnStats.PassTds = espnSsnStats.Stats.PassTds
					ssnStats.PassInts = espnSsnStats.Stats.PassInts
				}
				out.SeasonStats = append(out.SeasonStats, ssnStats)
			}
		}
	} else {
		out = &Player{
			Id:        m.Harris.Id,
			FirstName: m.Harris.FirstName,
			LastName:  m.Harris.LastName,
			Name:      m.Harris.Name,
			Position:  Position(m.Harris.Position),
			Team:      m.Harris.Team,
		}
	}
	if m.Harris != nil {
		out.CustomStdRank = m.Harris.StdRank
		out.CustomPprRank = m.Harris.PPRRank
	}

	return
}
