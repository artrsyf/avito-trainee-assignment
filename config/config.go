package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Auth AuthConfig `mapstructure:"auth"`
}

type AuthConfig struct {
	AccessTokenExpiration  string `mapstructure:"access_token_expiration"`
	RefreshTokenExpiration string `mapstructure:"refresh_token_expiration"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

func (c *AuthConfig) GetAccessTokenExpiration() (time.Duration, error) {
	return time.ParseDuration(c.AccessTokenExpiration)
}

func (c *AuthConfig) GetRefreshTokenExpiration() (time.Duration, error) {
	return time.ParseDuration(c.RefreshTokenExpiration)
}
