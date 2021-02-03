package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

type Post struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Slug        string    `json:"slug" db:"slug"`
	Content     string    `json:"content" db:"content"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	PublishedAt null.Time `json:"published_at" db:"published_at"`
	DeletedAt   null.Time `json:"-" db:"deleted_at"`
}
