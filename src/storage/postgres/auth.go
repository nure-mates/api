package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nure-mates/api/src/models"
)

type AuthRepo struct {
	*Postgres
}

func (p *Postgres) NewAuthRepo() *AuthRepo {
	return &AuthRepo{p}
}

func (ar *AuthRepo) CreateSession(ctx context.Context, session *models.UserSession) error {
	_, err := ar.DB.NewInsert().
		Model(session).
		Exec(ctx)

	return err
}

func (ar *AuthRepo) DisableSessionByID(ctx context.Context, sessionID uuid.UUID) error {
	expiredAt := time.Now().UTC()
	session := models.UserSession{
		ID:        sessionID,
		ExpiredAt: &expiredAt,
	}
	_, err := ar.DB.NewUpdate().
		Model(&session).
		WherePK().
		OmitZero().
		Exec(ctx)

	return err
}

func (ar *AuthRepo) GetSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (*models.UserSession, error) {
	session := &models.UserSession{}
	err := ar.DB.NewSelect().Model(session).
		Where("token_id = ?", tokenID).
		Limit(1).
		Scan(ctx)

	return session, err
}

func (ar *AuthRepo) UpdateSession(ctx context.Context, userSession *models.UserSession) error {
	_, err := ar.DB.NewUpdate().
		Model(userSession).
		WherePK().
		Exec(ctx)

	return err
}
