//go:generate mockgen -source=provider.go -destination=mocks/mock.go
package database

import (
	"context"
)

type DatabaseInterface interface {
	Exec(context.Context, string, ...interface{}) error
	Get(context.Context, interface{}, string, ...interface{}) error
	Select(context.Context, interface{}, string, ...interface{}) error
	QueryRow(context.Context, interface{}, string, ...interface{}) error
}

type DatabaseTx interface {
	DatabaseInterface
	Rollback() error
	Commit() error
}

type DatabasePrivoder interface {
	DatabaseInterface
	BeginTx(context.Context) (DatabaseTx, error)
}
