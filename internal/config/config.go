package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	RateLimiter RateLimiterConfig `mapstructure:"rate_limiter"`
}

type RateLimiterConfig struct {
	Requests int           `mapstructure:"requests"`
	Duration time.Duration `mapstructure:"duration"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	log.Printf("Конфиг: Requests = %d, Duration = %v", config.RateLimiter.Requests, config.RateLimiter.Duration)
	return config, nil
}
