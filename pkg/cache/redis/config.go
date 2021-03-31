package redis

import (
	"fmt"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database int
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d",
		c.Host, c.Port,
	)
}
