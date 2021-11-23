package models

type Room struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	HostID    int    `json:"host_id"`
	UserCount int    `json:"user_count" bun:"-"`
	UserIDs   []int  `bun:"m2m:user_to_room,join:Room=User"`
}

type UserToRoom struct {
	ID     int `bun:",pk"`
	UserID int `bun:"rel:belongs-to,join:user_id=id"`
	RoomID int `bun:"rel:belongs-to,join:room_id=id"`
}
