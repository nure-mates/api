package postgres

import (
	"context"
	"github.com/nure-mates/api/src/models"
)

type ProfileRepo struct {
	*Postgres
}

func (p *Postgres) NewProfileRepo() *ProfileRepo {
	return &ProfileRepo{p}
}

func (p *ProfileRepo) GetProfilesByEmail(ctx context.Context, email string) ([]models.User, error) {
	var res []models.User

	err := p.DB.NewSelect().
		Model(&res).
		Where("?TableAlias.email = ?", email).
		Scan(ctx)
	return res, toServiceError(err)
}
