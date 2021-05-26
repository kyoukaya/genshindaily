package models

type CheckInStatus int

const (
	CheckInStatusOK CheckInStatus = 1 + iota
	CheckInStatusDupe
	CheckInStatusFirstBind
)

func (c CheckInStatus) String() string {
	switch c {
	case CheckInStatusOK:
		return "Checked in successfully"
	case CheckInStatusDupe:
		return "Already checked in"
	case CheckInStatusFirstBind:
		return "First bind required"
	}
	return "UNKNOWN"
}

type SignInRequest struct {
	ActID string `json:"act_id"`
}

type OuterResponse struct {
	Retcode int64       `json:"retcode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type InfoResponse struct {
	TotalSignDay int64  `json:"total_sign_day"`
	Today        string `json:"today"`
	IsSign       bool   `json:"is_sign"`
	FirstBind    bool   `json:"first_bind"`
}

type RewardsResponse struct {
	Month  int64   `json:"month"`
	Awards []Award `json:"awards"`
}

type Award struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
	Cnt  int64  `json:"cnt"`
}
