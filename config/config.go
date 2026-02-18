package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Path       string
	ServerPort string `mapstructure:"SERVER_PORT"`
	DBConfig   `mapstructure:",squash"`
}

type DBConfig struct {
	DBUrl        string `mapstructure:"DB_URL"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	MigrationRun bool   `mapstructure:"MIGRATION_RUN"`
}

func NewConfig(path string) *Config {
	return &Config{Path: path}
}

func (c *Config) Load() error {
	v := viper.New()
	v.AddConfigPath(c.Path)
	v.SetConfigName("grpc-orders")
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("v.ReadInConfig: %w", err)
	}

	if err := v.Unmarshal(&c); err != nil {
		return fmt.Errorf("v.Unmarshal: %w", err)
	}

	return nil
}
