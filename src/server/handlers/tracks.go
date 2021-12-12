package handlers

import (
	"net/http"

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
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	//userID := r.Header.Get(models.UserIDHeaderName)
	//id, err := strconv.Atoi(userID)
	//if err != nil {
	//	SendEmptyResponse(w, http.StatusBadRequest)
	//	return
	//}

	track := models.Track{
		TrackURL: req.TrackURL,
		AddedBy:  10,
	}

	if err := h.service.AddTrack(r.Context(), &track); err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)
		return
	}

	SendEmptyResponse(w, http.StatusCreated)
}
