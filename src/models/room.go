package models

type Room struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	HostID    int    `json:"host_id"`
	UserCount int    `json:"user_count" bun:"-"`
	Users   []User `bun:"m2m:users_rooms,join:User=Room"`
}

type UsersRooms struct {
	ID     int `bun:",pk"`
	UserID int `json:"user_id"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`
	RoomID int `json:"room_id"`
	Room   *Room `bun:"rel:belongs-to,join:room_id=id"`
}
