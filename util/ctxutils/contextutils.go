package ctxutils

import (
	"context"

	"github.com/google/uuid"
)

// Custom type for context keys to avoid collisions when using context.WithValue
type ctxKey string

const userIDKey ctxKey = "userID"

// GetUserIDFromContext retrieves the user ID from the context.
// It returns the user ID and a boolean indicating whether it was found.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}

// SetUserIDInContext stores a UUID userID in the context.
func SetUserIDInContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
