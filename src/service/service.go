package service

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/nure-mates/api/src/config"
	"github.com/nure-mates/api/src/models"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	cfg         *config.Config
	authRepo    AuthRepo
	profileRepo ProfileRepo
	roomRepo    RoomRepo
}

func New(
	cfg *config.Config,
	aur AuthRepo,
	pr ProfileRepo,
	rr RoomRepo,
) *Service {
	once.Do(func() {
		service = &Service{
			cfg:         cfg,
			authRepo:    aur,
			profileRepo: pr,
			roomRepo:    rr,
		}
	})

	return service
}

func Get() *Service {
	return service
}

type AuthRepo interface {
	CreateSession(ctx context.Context, session *models.UserSession) error
	DisableSessionByID(ctx context.Context, sessionID uuid.UUID) error
	GetSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (*models.UserSession, error)
	UpdateSession(ctx context.Context, userSession *models.UserSession) error
}

type ProfileRepo interface {
	GetProfilesByEmail(ctx context.Context, email string) (models.User, error)
	AddNewProfile(ctx context.Context, newUser models.User) (models.User, error)
}

type RoomRepo interface {
	CreateRoom(ctx context.Context, room *models.Room) error
	GetUserRooms(ctx context.Context, userID int) ([]models.Room, error)
	GetRoom(ctx context.Context, id int) (*models.Room, error)
	AddUserToRoom(ctx context.Context, roomID, userID int) error
	RemoveUserFromRoom(ctx context.Context, roomID, userID int) error
	DeleteRoom(ctx context.Context, id int) error
	GetAvailableRooms(ctx context.Context, userID int) ([]models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	CheckRoom(ctx context.Context, id int) (bool, error)
}
