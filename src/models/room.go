package models

const (
	TokenLenForRoom = 16
)

// swagger:model
type Room struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	HostID int    `json:"host_id"`
	Public bool   `json:"public"`
	Token  string `json:"invite_token"`

	UserCount int     `json:"user_count" bun:"-"`
	Users     []User  `bun:"m2m:users_rooms,join:User=Room"`
	Tracks    []Track `bun:"rel:has-many,join:id=room_id"`
}

type UsersRooms struct {
	ID     int `bun:",pk"`
	UserID int `json:"user_id"`
	RoomID int `json:"room_id"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
	Room *Room `bun:"rel:belongs-to,join:room_id=id"`
}
