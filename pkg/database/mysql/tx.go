package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type MySQLTx struct {
	tx *sqlx.Tx
}

func (t *MySQLTx) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.tx.ExecContext(ctx, query, args...)
	return err
}

func (t *MySQLTx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.GetContext(ctx, dest, query, args...)
}

func (t *MySQLTx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.SelectContext(ctx, dest, query, args...)
}

func (t *MySQLTx) QueryRow(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	result := t.tx.QueryRowContext(ctx, query, args...)
	return result.Scan(dest)
}

func (t *MySQLTx) Rollback() error {
	return t.tx.Rollback()
}

func (t *MySQLTx) Commit() error {
	return t.tx.Commit()
}
