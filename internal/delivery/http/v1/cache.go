package v1

import (
	"context"
	"encoding/json"
	"log"
)

func (h *Handler) getFromCache(ctx context.Context, key string, out interface{}) (bool, error) {
	value, err := h.cache.Get(ctx, key)
	if err != nil {
		log.Printf("%s: get by key - %s", err, key)
		return false, nil
	}

	if err := json.Unmarshal(value, out); err != nil {
		return true, err
	}

	return true, nil
}

func (h *Handler) saveToCache(ctx context.Context, key string, input interface{}) {
	value, err := json.Marshal(input)
	if err == nil {
		if err := h.cache.Set(ctx, key, value); err != nil {
			log.Printf("%s: set by key -  %s", err, key)
		}
	}
}
