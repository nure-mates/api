package service

import (
	"context"

	"github.com/nure-mates/api/src/models"
	log "github.com/sirupsen/logrus"
)

func (s *Service) AddTrack(ctx context.Context, track *models.Track) error {
	if err := s.trackRepo.AddTrack(ctx, track); err != nil {
		log.Errorf("add track %s: %v", track.TrackURL, err)
		return err
	}
	return nil
}
