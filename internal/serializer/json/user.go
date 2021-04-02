package json

import (
	"encoding/json"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
)

type JsonUserSerializer struct{}

func NewUserSerializer() *JsonUserSerializer {
	return new(JsonUserSerializer)
}

func (s *JsonUserSerializer) Serialize(user domain.User) ([]byte, error) {
	return json.Marshal(user)
}

func (s *JsonUserSerializer) Deserialize(value []byte) (domain.User, error) {
	var user domain.User
	err := json.Unmarshal(value, &user)
	return user, err
}
