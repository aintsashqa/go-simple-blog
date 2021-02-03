package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      AppConfig
		Database MySQLDatabaseConfig
		Auth     AuthorizationConfig
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
		DatabaseName string `mapstructure:"db_name"`
		Charset      string `mapstructure:"charset"`
	}

	AuthorizationConfig struct {
		JWTSigningKey  string        `mapstructure:"jwt_signing_key"`
		JWTExpiresTime time.Duration `mapstructure:"jwt_expires_time"`
	}
)

var (
	prefixVariablesMap = map[string][]string{
		"app":  {"host", "port"},
		"db":   {"host", "port", "username", "password", "db_name", "charset"},
		"auth": {"jwt_signing_key", "jwt_expires_time"},
	}
)

func Init(filename string) (Config, error) {
	var cfg Config

	if err := parseConfigFile(filename); err != nil {
		return cfg, err
	}
	if err := parseEnv(); err != nil {
		return cfg, err
	}

	if err := unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func parseConfigFile(filename string) error {
	viper.AddConfigPath("configs")
	viper.SetConfigName(filename)
	return viper.ReadInConfig()
}

func parseEnvVariables(prefix string, variables ...string) error {
	viper.SetEnvPrefix(prefix)
	for _, v := range variables {
		if err := viper.BindEnv(v); err != nil {
			return err
		}
	}
	return nil
}

func parseEnv() error {
	for prefix, variables := range prefixVariablesMap {
		if err := parseEnvVariables(prefix, variables...); err != nil {
			return err
		}
	}

	return nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("app", &cfg.App); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("db", &cfg.Database); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}

	return nil
}
