package domain

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

const (
	CreateUserValidationAction UserValidationAction = iota
	UpdateUserValidationAction

	CreatePostValidationAction PostValidationAction = iota
	UpdatePostValidationAction
)

var (
	// User model errors
	ErrUserEmailEmptyValue       error = errors.New("Field email is required.")
	ErrUserEmailInvalidLength    error = errors.New("Field email must be greater than 6 and less 255 characters.")
	ErrUserEmailInvalidValue     error = errors.New("Field email must be an email.")
	ErrUserUsernameEmptyValue    error = errors.New("Field username is required.")
	ErrUserUsernameInvalidLength error = errors.New("Field username must be greater than 3 and less 255 characters.")
	ErrUserPasswordEmptyValue    error = errors.New("Field password is required.")
	ErrUserPasswordInvalidLength error = errors.New("Field password must be greater than 5 and less 255 characters.")

	// Post model errors
	ErrPostTitleEmptyValue      error = errors.New("Field title is required.")
	ErrPostTitleInvalidLength   error = errors.New("Field title must be greater than 8 and less 255 characters.")
	ErrPostSlugInvalidLength    error = errors.New("Field slug must be greater than 8 and less 255 characters.")
	ErrPostContentEmptyValue    error = errors.New("Field content is required.")
	ErrPostContentInvalidLength error = errors.New("Field content must be greater than 500 characters.")
)

type (
	UserValidationAction uint8
	PostValidationAction uint8

	Model struct {
		ID        uuid.UUID `json:"id"            db:"id"`
		CreatedAt time.Time `json:"created_at"    db:"created_at"`
		UpdatedAt time.Time `json:"updated_at"    db:"updated_at"`
		DeletedAt null.Time `json:"-"             db:"deleted_at"`
	}

	User struct {
		Model
		Email    string `json:"email,omitempty"    db:"email"`
		Username string `json:"username"           db:"username"`
		Password string `json:"-"                  db:"encrypted_password"`
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

func (u *User) Validate(action UserValidationAction) error {
	switch action {

	case CreateUserValidationAction:
		// Email validations
		if err := validation.Validate(&u.Email, validation.Required); err != nil {
			return ErrUserEmailEmptyValue
		}
		if err := validation.Validate(&u.Email, validation.Length(6, 255)); err != nil {
			return ErrUserEmailInvalidLength
		}
		if err := validation.Validate(&u.Email, is.Email); err != nil {
			return ErrUserEmailInvalidValue
		}

		// Password validations
		if err := validation.Validate(&u.Password, validation.Required); err != nil {
			return ErrUserPasswordEmptyValue
		}
		if err := validation.Validate(&u.Password, validation.Length(5, 255)); err != nil {
			return ErrUserPasswordInvalidLength
		}

	case UpdateUserValidationAction:
		// Username validations
		if err := validation.Validate(&u.Username, validation.Required); err != nil {
			return ErrUserUsernameEmptyValue
		}
		if err := validation.Validate(&u.Username, validation.Length(3, 255)); err != nil {
			return ErrUserUsernameInvalidLength
		}

	}

	return nil
}

func (p *Post) Validate(action PostValidationAction) error {
	switch action {

	case CreatePostValidationAction, UpdatePostValidationAction:
		// Title validations
		if err := validation.Validate(&p.Title, validation.Required); err != nil {
			return ErrPostTitleEmptyValue
		}
		if err := validation.Validate(&p.Title, validation.Length(8, 255)); err != nil {
			return ErrPostTitleInvalidLength
		}

		// Slug validations
		if err := validation.Validate(&p.Slug, validation.Length(8, 255)); err != nil {
			return ErrPostSlugInvalidLength
		}

		// Content validations
		if err := validation.Validate(&p.Content, validation.Required); err != nil {
			return ErrPostContentEmptyValue
		}
		if err := validation.Validate(&p.Content, validation.Length(500, 0)); err != nil {
			return ErrPostContentInvalidLength
		}

	}

	return nil
}
