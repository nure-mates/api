package service

import (
	"crypto/subtle"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"

	"github.com/nure-mates/api/src/models"
	"github.com/nure-mates/api/src/storage/redis"
)

// Token error list.
var (
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrTokenInBlackList   = errors.New("token in black list")
	ErrTokenClaimsInvalid = errors.New("token claims are invalid")
)

// Validate validates access token.
func (s *Service) Validate(accessToken string) (*models.UserSession, error) {
	token, err := s.parseJWT(accessToken)
	if err != nil {
		return nil, err
	}

	claims := s.parseClaims(token)
	if claims == nil || !token.Valid {
		return nil, ErrTokenInvalid
	}

	if err := s.checkBlackList(claims.TokenID.String()); err != nil {
		return nil, err
	}

	return &models.UserSession{
		UserID:  claims.UserID,
		TokenID: claims.TokenID,
	}, nil
}

// ValidateExternalAPIToken validates token used for service to service comunication.
func (s *Service) ValidateExternalAPIToken(token string) error {
	if subtle.ConstantTimeCompare([]byte(token), []byte(s.cfg.ExternalAPIToken)) != 1 {
		return ErrTokenInvalid
	}

	return nil
}

func (s *Service) checkBlackList(tokenID string) error {
	res, err := s.redis.Get(redis.TokenBlackListKey(tokenID))
	if err != nil {
		return err
	}

	if res != nil {
		return ErrTokenInBlackList
	}

	return nil
}

func (s *Service) parseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.AccessTokenSecret), nil
	})

	if token == nil {
		return nil, ErrTokenInvalid
	}

	//nolint: errorlint
	switch ve := err.(type) {
	case *jwt.ValidationError:
		if ve.Errors|(jwt.ValidationErrorExpired) != jwt.ValidationErrorExpired {
			return nil, ErrTokenInvalid
		}
	case nil:
	default:
		return nil, err
	}

	return token, nil
}
