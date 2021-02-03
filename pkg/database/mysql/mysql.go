package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func NewMySQL(cfg Config) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", cfg.Dsn())
}

func Migrate(connection *sqlx.DB) error {
	driver, err := mysql.WithInstance(connection.DB, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"mysql", driver,
	)

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}
