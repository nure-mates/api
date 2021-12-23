package models

const (
	TokenLenForRoom = 16
)

// swagger:model
type Room struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	HostID    int     `json:"host_id"`
	Auto      bool    `json:"auto"`
	Fixed     int     `json:"fixed_count"`
	Public    *bool   `json:"public"`
	Token     string  `json:"invite_token"`
	MaxUsers  int     `json:"max_users"`
	UserCount int     `json:"user_count" bun:"-"`
	Users     []User  `bun:"m2m:users_rooms,join:User=Room"`
	Tracks    []Track `bun:"rel:has-many,join:id=room_id"`
}

type UsersRooms struct {
	ID     int   `bun:",pk"`
	UserID int   `json:"user_id"`
	RoomID int   `json:"room_id"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`
	Room   *Room `bun:"rel:belongs-to,join:room_id=id"`
}
