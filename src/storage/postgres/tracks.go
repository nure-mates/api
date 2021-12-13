package postgres

import (
	"context"

	"github.com/nure-mates/api/src/models"
)

type TrackRepo struct {
	*Postgres
}

func (p *Postgres) NewTrackRepo() *TrackRepo {
	p.DB.RegisterModel((*models.Track)(nil))
	return &TrackRepo{p}
}

func (r *TrackRepo) AddTrack(ctx context.Context, track *models.Track) error {
	_, err := r.DB.NewInsert().
		Model(track).
		Returning("id").
		Exec(ctx)

	return err
}
