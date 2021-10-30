package service

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/nure-mates/api/src/config"
	"github.com/nure-mates/api/src/models"
	"github.com/nure-mates/api/src/storage/redis"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	cfg         *config.Config
	redis       *redis.Client
	authRepo    AuthRepo
	profileRepo ProfileRepo
}

func New(
	cfg *config.Config,
	rds *redis.Client,
	aur AuthRepo,
	pr ProfileRepo,
) *Service {
	once.Do(func() {
		service = &Service{
			redis:       rds,
			cfg:         cfg,
			authRepo:    aur,
			profileRepo: pr,
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
	GetProfilesByEmail(ctx context.Context, email string) ([]models.User, error)
}