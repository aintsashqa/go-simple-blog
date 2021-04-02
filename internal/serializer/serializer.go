package serializer

import (
	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/serializer/json"
)

type (
	UserSerializer interface {
		Serialize(domain.User) ([]byte, error)
		Deserialize([]byte) (domain.User, error)
	}

	PostSerializer interface {
		Serialize(domain.Post) ([]byte, error)
		Deserialize([]byte) (domain.Post, error)
	}

	Serializer struct {
		User UserSerializer
		Post PostSerializer
	}
)

func NewSerializer() *Serializer {
	return &Serializer{
		User: json.NewUserSerializer(),
		Post: json.NewPostSerializer(),
	}
}
