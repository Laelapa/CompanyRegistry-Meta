package domain

import (
	"github.com/google/uuid"
)

type User struct {
	// Fields kept as pointers for less friction if implementing
	// partial updates is decided in the future
	ID           *uuid.UUID
	Username     *string
	PasswordHash *string
}
