package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Username  string    `db:"username"`
	Password  string    `db:"encrypted_password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt null.Time `db:"deleted_at"`
}
