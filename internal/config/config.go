package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      AppConfig           `mapstructure:"app"`
		Database MySQLDatabaseConfig `mapstructure:"db"`
		Auth     AuthorizationConfig `mapstructure:"auth"`
	}

	AppConfig struct {
		Host            string        `mapstructure:"host"`
		Port            int           `mapstructure:"port"`
		ReadTimeout     time.Duration `mapstructure:"read_timeout"`
		WriteTimeout    time.Duration `mapstructure:"write_timeout"`
		MaxHeaderMBytes int           `mapstructure:"max_header_mbytes"`
	}

	MySQLDatabaseConfig struct {
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		Username     string `mapstructure:"username"`
		Password     string `mapstructure:"password"`
		DatabaseName string `mapstructure:"name"`
		Charset      string `mapstructure:"charset"`
	}

	AuthorizationConfig struct {
		JWTSigningKey  string        `mapstructure:"jwt_signing_key"`
		JWTExpiresTime time.Duration `mapstructure:"jwt_expires_time"`
	}
)

func Init(filename string) (Config, error) {
	var cfg Config

	viper.AddConfigPath("configs")
	viper.SetConfigName(filename)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	err := viper.Unmarshal(&cfg)
	return cfg, err
}
