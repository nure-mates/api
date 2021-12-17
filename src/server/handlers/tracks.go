package handlers

import (
	"net/http"

	"github.com/nure-mates/api/src/context"
	log "github.com/sirupsen/logrus"

	"github.com/nure-mates/api/src/models"
	"github.com/nure-mates/api/src/service"
)

type TrackHandler struct {
	service *service.Service
}

func NewTrackHandler(s *service.Service) *TrackHandler {
	return &TrackHandler{
		service: s,
	}
}

func (h *TrackHandler) AddTrack(w http.ResponseWriter, r *http.Request) {
	req := &models.AddTrackRequest{}

	err := UnmarshalRequest(r, req)
	if err != nil {
		log.Errorf("failed to unmarshal request: %v", err)
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	id := context.GetUserID(r.Context())

	track := models.Track{
		TrackURL: req.TrackURL,
		AddedBy:  id,
		RoomID:   req.RoomID,
	}

	log.Printf("User is %d\n", id)

	if err = h.service.AddTrack(r.Context(), &track); err != nil {
		log.Errorf("failed to add track: %v", err)
		SendResponse(w, http.StatusBadRequest, err)
		return
	}

	SendEmptyResponse(w, http.StatusCreated)
}
