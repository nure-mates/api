package service

import (
	"context"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
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
	spotifyAuth *spotifyauth.Authenticator
}

func New(
	cfg *config.Config,
	aur AuthRepo,
	pr ProfileRepo,
	authenticator *spotifyauth.Authenticator,
) *Service {
	once.Do(func() {
		service = &Service{
			cfg:         cfg,
			authRepo:    aur,
			profileRepo: pr,
			spotifyAuth: authenticator,
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
