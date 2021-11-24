package handlers

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/nure-mates/api/src/models"
	"github.com/nure-mates/api/src/service"
)

type AuthHandler struct {
	service *service.Service
}

func NewAuthHandler(s *service.Service) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

// swagger:operation POST /login auth login
//   create a session and obtain tokens pair
// ---
// parameters:
// - name: login_request
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/LoginRequest'
// responses:
//   '200':
//     description: Fetched
//     schema:
//       "$ref": "#/definitions/LoginResponse"
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := &models.LoginRequest{}

	err := UnmarshalRequest(r, req)
	if err != nil {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	res, err := h.service.Login(r.Context(), *req)
	if err != nil {
		SendEmptyResponse(w, http.StatusInternalServerError)
		return
	}

	SendResponse(w, http.StatusOK, res)
}

// SpotifyLogin description should be written later
func (h *AuthHandler) SpotifyLogin(w http.ResponseWriter, r *http.Request) {
	authURL := h.service.GetSpotifyAuthUrl(service.MockState)

	http.Redirect(w, r, authURL, http.StatusMovedPermanently)
}

// SpotifyCallback description should be written later
func (h *AuthHandler) SpotifyCallback(w http.ResponseWriter, r *http.Request) {
	email, err := h.service.GetSpotifyData(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusIMUsed)
	}

	respJson := fmt.Sprintf(`{"email": "%s"}`, email)

	SendResponse(w, http.StatusOK, []byte(respJson))
}

// swagger:operation DELETE /logout auth logout
//   deactivate user session, move access token to the black list
// ---
// responses:
//   '204':
//     description: Successfully logged out
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get(AccessTokenHeader)

	jwtAccessToken, err := ParseAuthorizationHeader(accessToken, BearerSchema)
	if err != nil {
		log.Error("logout error ", err)
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	err = h.service.Logout(r.Context(), jwtAccessToken)
	if err != nil {
		log.Error("logout error ", err)
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendResponse(w, http.StatusNoContent, nil)
}

// swagger:operation POST /token auth token
//   refresh access token if previous tokens pair was valid
// ---
// parameters:
// - name: TokenPair
//   in: body
//   required: true
//   schema:
//     $ref: '#/definitions/TokenPair'
// responses:
//   '201':
//     description: Created
//     schema:
//       "$ref": "#/definitions/TokenPair"
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (h *AuthHandler) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	req := &models.TokenPair{}

	err := UnmarshalRequest(r, req)
	if err != nil {
		log.Error("token refresh error ", err)
		SendEmptyResponse(w, http.StatusBadRequest)

		return
	}

	res, err := h.service.RefreshToken(r.Context(), req)
	if err != nil {
		log.Error("token refresh error ", err)
		SendEmptyResponse(w, http.StatusInternalServerError)

		return
	}

	SendResponse(w, http.StatusCreated, res)
}
