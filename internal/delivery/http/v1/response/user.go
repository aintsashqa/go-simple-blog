package response

import (
	"time"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	uuid "github.com/satori/go.uuid"
)

type TokenResponseDto struct {
	AccessToken string `json:"access_token"`
}

func (dto *TokenResponseDto) TransformFromObject(tokens service.Tokens) {
	dto.AccessToken = tokens.AccessToken
}

type UserResponseDto struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (dto *UserResponseDto) TransformFromObject(user domain.User) {
	dto.ID = user.ID
	dto.Username = user.Username
	dto.CreatedAt = user.CreatedAt
	dto.UpdatedAt = user.UpdatedAt
}
