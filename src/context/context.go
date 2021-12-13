package context

import (
	"context"
)

const (
	userIDKey = "device_id"
)

// GetUserID - returns User UID from context.
func GetUserID(ctx context.Context) int {
	value, _ := ctx.Value(userIDKey).(int)

	return value
}

// WithUserID - add User UID value to context.
func WithUserID(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, userIDKey, value) //nolint
}
