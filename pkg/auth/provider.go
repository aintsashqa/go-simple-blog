//go:generate mockgen -source=provider.go -destination=mocks/mock.go
package auth

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TokenParams struct {
	UserID    uuid.UUID
	ExpiresAt time.Duration
}

type AuthorizationProvider interface {
	NewToken(TokenParams) (string, error)
	Parse(string) (uuid.UUID, error)
}
