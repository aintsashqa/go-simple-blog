//go:generate mockgen -source=provider.go -destination=mocks/mock.go
package database

import (
	"context"
)

type DatabasePrivoder interface {
	Exec(context.Context, string, ...interface{}) error
	Get(context.Context, interface{}, string, ...interface{}) error
	Select(context.Context, interface{}, string, ...interface{}) error
	QueryRow(context.Context, interface{}, string, ...interface{}) error
}
