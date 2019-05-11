package db

type UserLocation struct {
	UserId string `json:"user_id"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
