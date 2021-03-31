package cache

import (
	"context"
)

type CachePrivoder interface {
	Set(context.Context, string, []byte) error
	Get(context.Context, string) ([]byte, error)
}
