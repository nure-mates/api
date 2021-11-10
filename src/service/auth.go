package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"

	"github.com/nure-mates/api/src/models"
)

const hoursInDay = 24

var (
	issMap = map[string]struct{}{
		"accounts.google.com":         {},
		"https://accounts.google.com": {},
	}

	ErrSessionNotFound  = errors.New("session by token ID was not found")
	ErrTokensMismatched = errors.New("access and refresh tokens are mismatched")
	ErrSessionExpired   = errors.New("session expired")
)

func (s *Service) Login(ctx context.Context, loginReq models.LoginRequest) (resp models.LoginResponse, err error) {
	tokenInfo, err := verifyIDToken(ctx, loginReq.IDToken, s.cfg.GoogleAud)
	if err != nil {
		log.Error("login error, validate ID token ", err)
		return resp, err
	}

	user, err := s.checkUser(ctx, tokenInfo.Claims["email"].(string))
	if err != nil {
		log.Error("login error, check user ", err)
		return resp, err
	}

	resp.User = user

	claims := models.Claims{
		SessionID: uuid.New(),
		TokenID:   uuid.New(),
		UserID:    user.ID,
	}

	accessToken, err := s.GenerateAccess(&claims)
	if err != nil {
		return resp, err
	}

	resp.TokenPair.AccessToken = accessToken

	refreshToken, err := s.GenerateRefresh()
	if err != nil {
		return resp, errors.Wrap(err, "generate refresh token")
	}

	resp.TokenPair.RefreshToken = refreshToken

	err = s.createSession(ctx, &claims, &resp.TokenPair)
	if err != nil {
		return resp, errors.Wrap(err, "can't create user session")
	}

	return resp, err
}

func (s *Service) Logout(ctx context.Context, accessToken string) (err error) {
	claims, err := s.Revoke(accessToken)
	if err != nil {
		return err
	}

	err = s.authRepo.DisableSessionByID(ctx, claims.SessionID)
	if err != nil {
		return err
	}

	return nil
}

// RefreshToken refreshes access token.
func (s *Service) RefreshToken(ctx context.Context, oldTokens *models.TokenPair) (*models.TokenPair, error) {
	access, err := s.parseJWT(oldTokens.AccessToken)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	accessClaims := s.parseClaims(access)
	if accessClaims == nil {
		return nil, ErrTokenClaimsInvalid
	}

	session, err := s.authRepo.GetSessionByTokenID(ctx, accessClaims.TokenID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.RefreshToken != oldTokens.RefreshToken {
		return nil, ErrTokensMismatched
	}

	now := time.Now().UTC()
	if now.After(*session.ExpiredAt) {
		return nil, ErrSessionExpired
	}

	claims := models.Claims{
		TokenID:   uuid.New(),
		SessionID: session.ID,
		UserID:    session.UserID,
	}

	accessNew, err := s.GenerateAccess(&claims)
	if err != nil {
		return nil, errors.Wrap(err, "token generate ")
	}

	session.TokenID = claims.TokenID
	session.UpdatedAt = &now

	err = s.authRepo.UpdateSession(ctx, session)
	if err != nil {
		return nil, errors.Wrap(err, "update user session ")
	}

	if _, err = s.Revoke(oldTokens.AccessToken); err != nil {
		return nil, err
	}

	return &models.TokenPair{AccessToken: accessNew, RefreshToken: oldTokens.RefreshToken}, nil
}

func verifyIDToken(ctx context.Context, idToken, aud string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(ctx, idToken, aud)

	if err != nil {
		log.Error("validate id token error ", err)
		return nil, err
	}

	if _, ok := issMap[payload.Issuer]; !ok {
		log.Error("validate id token error, iss is not valid ", err)
		return nil, err
	}

	emailVerified, ok := payload.Claims["email_verified"].(bool)
	if !ok || !emailVerified {
		log.Error("validate id token error, email is not verified ", err)
		return nil, err
	}

	return payload, nil
}

func (s *Service) checkUser(ctx context.Context, email string) (user models.User, err error) {
	userCandidates, err := s.profileRepo.GetProfilesByEmail(ctx, email)
	if err != nil {
		return user, err
	}

	for i := range userCandidates {
		if !userCandidates[i].Archived {
			return userCandidates[i], nil
		}
	}

	user, err = s.profileRepo.AddNewProfile(ctx, models.User{
		Email:     email,
		Archived:  false,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return user, err
	}

	return user, nil
}

// GenerateAccess generates token with claims.
func (s *Service) GenerateAccess(claims *models.Claims) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Unix() + int64(s.cfg.AccessTokenTTL)
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenWithClaims.SignedString([]byte(s.cfg.AccessTokenSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// GenerateRefresh generates refresh token.
func (s *Service) GenerateRefresh() (string, error) {
	return generateRandomString(s.cfg.RefreshTokenLen)
}

func (s *Service) parseClaims(token *jwt.Token) *models.Claims {
	if claims, ok := token.Claims.(*models.Claims); ok {
		return claims
	}

	return nil
}

// getExpiredAt returns session expiration time till the end of the current day.
func getExpiredAt() time.Time {
	return time.Now().UTC().AddDate(0, 0, 1).Truncate(time.Hour * hoursInDay)
}

func (s *Service) createSession(ctx context.Context, claims *models.Claims, tokenPair *models.TokenPair) error {
	now := time.Now().UTC()
	expiredAt := getExpiredAt()
	session := models.UserSession{
		ID:           claims.SessionID,
		UserID:       claims.UserID,
		TokenID:      claims.TokenID,
		RefreshToken: tokenPair.RefreshToken,
		CreatedAt:    &now,
		UpdatedAt:    &now,
		ExpiredAt:    &expiredAt,
	}

	return s.authRepo.CreateSession(ctx, &session)
}

// Revoke revokes access token.
func (s *Service) Revoke(accessToken string) (*models.Claims, error) {
	token, err := s.parseJWT(accessToken)
	if err != nil {
		return nil, err
	}

	claims := s.parseClaims(token)
	if claims == nil {
		return nil, ErrTokenClaimsInvalid
	}

	return claims, nil
}

// generateRandomString returns a URL-safe, base64 encoded securely generated random string.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
