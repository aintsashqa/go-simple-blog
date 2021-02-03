package mysql

import (
	"fmt"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	Charset  string
}

func (c *Config) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset,
	)
}
