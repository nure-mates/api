package handlers

import (
	"context"

	"github.com/nure-mates/api/src/models"
)

type AuthService interface {
	Login(ctx context.Context, loginReq models.LoginRequest) (resp models.LoginResponse, err error)
	Logout(ctx context.Context, accessToken string) (err error)
	RefreshToken(ctx context.Context, tokenReq models.TokenPair) (resp models.TokenPair, err error)
}

type RoomService interface {
	CreateRoom(ctx context.Context, room *models.Room) error
	GetUserRooms(ctx context.Context, userID int) ([]models.Room, error)
	GetRoom(ctx context.Context, id int) *models.Room
	AddUserToRoom(ctx context.Context, userID int) error
	RemoveUserFromRoom(ctx context.Context, userID int) error
	DeleteRoom(ctx context.Context, id int) error
	GetAvailableRooms(ctx context.Context, userID int) ([]models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	CheckRoom(ctx context.Context, id int) (bool, error)
}
