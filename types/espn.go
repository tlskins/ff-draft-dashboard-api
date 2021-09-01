package types

import "strconv"

type EspnPlayersResp struct {
	Players []*EspnPlayer `json:"players"`
}

type ESPNPosition int

const (
	EspnQB ESPNPosition = 1
	EspnRB ESPNPosition = 2
	EspnWR ESPNPosition = 3
	EspnTE ESPNPosition = 4
)

type ESPNTeam int

const (
	EspnATL ESPNTeam = 1
	EspnBUF ESPNTeam = 2
	EspnCHI ESPNTeam = 3
	EspnCIN ESPNTeam = 4
	EspnCLE ESPNTeam = 5
	EspnDAL ESPNTeam = 6
	EspnDEN ESPNTeam = 7
	EspnDET ESPNTeam = 8
	EspnGB  ESPNTeam = 9
	EspnTEN ESPNTeam = 10
	EspnIND ESPNTeam = 11
	EspnKC  ESPNTeam = 12
	EspnLV  ESPNTeam = 13
	EspnLAR ESPNTeam = 14
	EspnMIA ESPNTeam = 15
	EspnMIN ESPNTeam = 16
	EspnNE  ESPNTeam = 17
	EspnNO  ESPNTeam = 18
	EspnNYG ESPNTeam = 19
	EspnNYJ ESPNTeam = 20
	EspnPHL ESPNTeam = 21
	EspnARI ESPNTeam = 22
	EspnPIT ESPNTeam = 23
	EspnLAC ESPNTeam = 24
	EspnSF  ESPNTeam = 25
	EspnSEA ESPNTeam = 26
	EspnTB  ESPNTeam = 27
	EspnWAS ESPNTeam = 28
	EspnCAR ESPNTeam = 29
	EspnJAC ESPNTeam = 30
	EspnBAL ESPNTeam = 33
	EspnHOU ESPNTeam = 34
)

type EspnPlayer struct {
	Id                int                `json:"id"`
	DraftAuctionValue int                `json:"draftAuctionValue"`
	OnTeamId          int                `json:"onTeamId"`
	Profile           *EspnPlayerProfile `json:"player"`
}

type EspnPlayerProfile struct {
	Id                int           `json:"id"`
	DefaultPositionId ESPNPosition  `json:"defaultPositionId"`
	EligibleSlots     []int         `json:"eligibleSlots"`
	FirstName         string        `json:"firstName"`
	LastName          string        `json:"lastName"`
	FullName          string        `json:"fullName"`
	Injured           bool          `json:"injured"`
	InjuryStatus      string        `json:"injuryStatus"`
	ProTeamId         ESPNTeam      `json:"proTeamId"`
	LastNewsDate      int           `json:"lastNewsDate"`
	SeasonOutlook     string        `json:"seasonOutlook"`
	Ranks             *ESPNRankings `json:"draftRanksByRankType"`
}

type ESPNRankings struct {
	PPR      *ESPNRank `json:"PPR"`
	Standard *ESPNRank `json:"STANDARD"`
}

type ESPNRank struct {
	AuctionValue int    `json:"auctionValue"`
	Rank         int    `json:"rank"`
	RankType     string `json:"rankType"`
}

func (p EspnPlayer) Position() Position {
	switch p.Profile.DefaultPositionId {
	case EspnQB:
		return QB
	case EspnRB:
		return RB
	case EspnWR:
		return WR
	case EspnTE:
		return TE
	default:
		return NoPosition
	}
}

func (p EspnPlayer) Team() (out string) {
	switch p.Profile.ProTeamId {
	case EspnATL:
		return "ATL"
	case EspnBUF:
		return "BUF"
	case EspnCHI:
		return "CHI"
	case EspnCIN:
		return "CIN"
	case EspnCLE:
		return "CLE"
	case EspnDAL:
		return "DAL"
	case EspnDEN:
		return "DEN"
	case EspnDET:
		return "DET"
	case EspnGB:
		return "GB"
	case EspnTEN:
		return "TEN"
	case EspnIND:
		return "IND"
	case EspnKC:
		return "KC"
	case EspnLV:
		return "LV"
	case EspnLAR:
		return "LAR"
	case EspnMIA:
		return "MIA"
	case EspnMIN:
		return "MIN"
	case EspnNE:
		return "NE"
	case EspnNO:
		return "NO"
	case EspnNYG:
		return "NYG"
	case EspnNYJ:
		return "NYJ"
	case EspnPHL:
		return "PHL"
	case EspnARI:
		return "ARI"
	case EspnPIT:
		return "PIT"
	case EspnLAC:
		return "LAC"
	case EspnSF:
		return "SF"
	case EspnSEA:
		return "SEA"
	case EspnTB:
		return "TB"
	case EspnWAS:
		return "WAS"
	case EspnCAR:
		return "CAR"
	case EspnJAC:
		return "JAC"
	case EspnBAL:
		return "BAL"
	case EspnHOU:
		return "HOU"
	default:
		return ""
	}
}

func (p EspnPlayer) ToPlayer() (out *Player) {
	out = &Player{
		Id:        strconv.Itoa(p.Id),
		FirstName: p.Profile.FirstName,
		LastName:  p.Profile.LastName,
		Name:      p.Profile.FullName,
		MatchName: MatchName(p.Profile.FullName),
		Position:  p.Position(),
		Team:      p.Team(),
	}
	if p.Profile.Ranks != nil && p.Profile.Ranks.PPR != nil {
		out.EspnOvrPprRank = p.Profile.Ranks.PPR.Rank
	}
	if p.Profile.Ranks != nil && p.Profile.Ranks.Standard != nil {
		out.EspnOvrStdRank = p.Profile.Ranks.Standard.Rank
	}

	return
}
