package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	UserIDHeaderName = "initiator-user-id"
	RoleIDHeaderName = "initiator-role-id"
)

type UserSession struct {
	ID           uuid.UUID  `json:"id"`
	TokenID      uuid.UUID  `json:"token_id"`
	UserID       int        `json:"user_id"`
	RefreshToken string     `json:"refresh_token"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	ExpiredAt    *time.Time `json:"expired_at"`
}

// Claims token's claims/payload.
type Claims struct {
	SessionID uuid.UUID
	TokenID   uuid.UUID
	UserID    int

	jwt.StandardClaims
}

// TTL returns TTL in seconds.
func (c *Claims) TTL() int64 {
	return c.StandardClaims.ExpiresAt - time.Now().Unix()
}
