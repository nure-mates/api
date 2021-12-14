package models

type AddTrackRequest struct {
	TrackURL string `json:"trackURL"`
}

type Track struct {
	ID       int
	TrackURL string
	AddedBy  int
	RoomID   int
}
