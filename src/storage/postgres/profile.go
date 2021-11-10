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

func (p *ProfileRepo) AddNewProfile(ctx context.Context, newUser models.User) (models.User, error) {

	_, err := p.DB.NewInsert().Model(&newUser).Exec(ctx)

	return newUser, toServiceError(err)
}
