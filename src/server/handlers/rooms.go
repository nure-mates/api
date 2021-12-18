package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/nure-mates/api/src/context"
	"github.com/nure-mates/api/src/models"
	"github.com/nure-mates/api/src/service"
)

type RoomHandler struct {
	service *service.Service
}

func NewRoomHandler(s *service.Service) *RoomHandler {
	return &RoomHandler{
		service: s,
	}
}

// swagger:operation POST /create-room rooms create_room
//   create a room
// ---
// parameters:
// - name: create_room_request
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/CreateRoomRequest'
// responses:
//   '201':
//     description: Created
//   '400':
//     description: Bad Request
//     schema:
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	room := &models.Room{}

	err := UnmarshalRequest(r, room)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	room, err = h.service.CreateRoom(r.Context(), room)
	if err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	w.Header().Add(headerForPrivateRoom, room.Token)

	SendEmptyResponse(w, http.StatusCreated)
}

func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["room-id"])
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	userID := context.GetUserID(r.Context())

	room, err := h.service.GetRoom(r.Context(), id, userID)
	if err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendResponse(w, http.StatusOK, map[string]interface{}{
		"room": room,
	})
}

func (h *RoomHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	room := &models.Room{}

	err := UnmarshalRequest(r, room)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	if err := h.service.UpdateRoom(r.Context(), room); err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendEmptyResponse(w, http.StatusOK)
}

//GetUserRooms ...
func (h *RoomHandler) GetUserRooms(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.Atoi(params["user-id"])
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	rooms, err := h.service.GetUserRooms(r.Context(), userID)
	if err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendResponse(w, http.StatusOK, rooms)
}

//GetPublicRooms ...
func (h *RoomHandler) GetPublicRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.GetPublicRooms(r.Context())
	if err != nil {
		SendHTTPError(w, err)
		SendEmptyResponse(w, http.StatusInternalServerError)
		return
	}

	SendResponse(w, http.StatusOK, rooms)
}

//GetAvailableRooms ...
func (h *RoomHandler) GetAvailableRooms(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.Atoi(params["user-id"])
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	rooms, err := h.service.GetAvailableRooms(r.Context(), userID)
	if err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendResponse(w, http.StatusOK, rooms)
}

//AddUserToRoom ...
func (h *RoomHandler) AddUserToRoom(w http.ResponseWriter, r *http.Request) {
	var (
		roomUserIDs models.UsersRooms
	)

	err := UnmarshalRequest(r, &roomUserIDs)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	if err := h.service.AddUserToRoom(r.Context(), roomUserIDs.RoomID, roomUserIDs.UserID); err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendEmptyResponse(w, http.StatusOK)
}

//RemoveUserFromRoom ...
func (h *RoomHandler) RemoveUserFromRoom(w http.ResponseWriter, r *http.Request) {
	var (
		roomUserIDs models.UsersRooms
	)

	err := UnmarshalRequest(r, &roomUserIDs)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	if err := h.service.RemoveUserFromRoom(r.Context(), roomUserIDs.RoomID, roomUserIDs.UserID); err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendEmptyResponse(w, http.StatusOK)
}

//DeleteRoom ...
func (h *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["room-id"])
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	if err := h.service.DeleteRoom(r.Context(), id); err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendEmptyResponse(w, http.StatusOK)
}

//JoinToRoom ...
func (h *RoomHandler) JoinToRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token := params["token"]

	userID := context.GetUserID(r.Context())

	err := h.service.JoinToRoom(r.Context(), token, userID)
	if err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)
		return
	}

	SendEmptyResponse(w, http.StatusOK)
}
