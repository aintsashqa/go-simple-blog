package json

import (
	"encoding/json"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
)

type JsonPostSerializer struct{}

func NewPostSerializer() *JsonPostSerializer {
	return new(JsonPostSerializer)
}

func (s *JsonPostSerializer) Serialize(post domain.Post) ([]byte, error) {
	return json.Marshal(post)
}

func (s *JsonPostSerializer) Deserialize(value []byte) (domain.Post, error) {
	var post domain.Post
	err := json.Unmarshal(value, &post)
	return post, err
}
