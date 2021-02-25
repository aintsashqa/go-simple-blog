package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

type (
	Model struct {
		ID        uuid.UUID `json:"id"            db:"id"`
		CreatedAt time.Time `json:"created_at"    db:"created_at"`
		UpdatedAt time.Time `json:"updated_at"    db:"updated_at"`
		DeletedAt null.Time `json:"-"             db:"deleted_at"`
	}

	User struct {
		Model
		Email    string `json:"email"        db:"email"`
		Username string `json:"username"     db:"username"`
		Password string `json:"-"            db:"encrypted_password"`
	}

	Post struct {
		Model
		Title       string    `json:"title"           db:"title"`
		Slug        string    `json:"slug"            db:"slug"`
		Content     string    `json:"content"         db:"content"`
		UserID      uuid.UUID `json:"user_id"         db:"user_id"`
		PublishedAt null.Time `json:"published_at"    db:"published_at"`
	}
)

func (m *Model) Init() {
	m.ID = uuid.NewV4()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.DeletedAt = null.NewTime(time.Now(), false)
}

func (m *Model) Update() {
	m.UpdatedAt = time.Now()
}

func (m *Model) Delete() {
	m.DeletedAt = null.NewTime(time.Now(), true)
}
