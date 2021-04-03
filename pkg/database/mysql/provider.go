package mysql

import (
	"context"

	"github.com/aintsashqa/go-simple-blog/pkg/database"
	"github.com/jmoiron/sqlx"
)

type MySQLProvider struct {
	db *sqlx.DB
}

func NewMySQLProvider(cfg Config) (*MySQLProvider, error) {
	db, err := NewMySQL(cfg)
	if err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}

	return &MySQLProvider{db: db}, nil
}

func (p *MySQLProvider) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *MySQLProvider) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.db.GetContext(ctx, dest, query, args...)
}

func (p *MySQLProvider) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.db.SelectContext(ctx, dest, query, args...)
}

func (p *MySQLProvider) QueryRow(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	result := p.db.QueryRowContext(ctx, query, args...)
	return result.Scan(dest)
}

func (p *MySQLProvider) BeginTx(ctx context.Context) (database.DatabaseTx, error) {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	result := &MySQLTx{tx: tx}
	return result, nil
}

func (p *MySQLProvider) Close() error {
	return p.db.Close()
}
