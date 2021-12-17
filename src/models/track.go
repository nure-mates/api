package models

type AddTrackRequest struct {
	TrackURL string `json:"trackURL"`
	RoomID   int    `json:"roomID"`
}

type Track struct {
	ID       int
	TrackURL string
	AddedBy  int
	RoomID   int
}
