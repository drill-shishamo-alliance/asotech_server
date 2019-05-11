package db

type RoomUser struct {
	UserId []string `json:"user_id"`
}

type RoomId struct {
	Value string `json:"room_id"`
}

type RemainingHuman struct {
	Value string `json:"remaining_human"`
}

type RoomMemberLimit struct {
	Value string `json:"room_member_limit"`
}