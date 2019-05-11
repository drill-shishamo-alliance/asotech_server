package view

type PostRoom struct {
	UserId string `json:"user_id"`
	RestrictTime int `json:"restrict_time"`
}

type PostBelongToTheRoom struct {
	UserId string `json:"user_id"`
}

type PostHumanCollaborate struct {
	UserId string `json:"user_id"`
}