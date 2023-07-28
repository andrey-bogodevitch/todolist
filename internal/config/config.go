package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr         []string `env:"REDIS_ADDR" required:"true"`
	RedisPassword     string   `env:"REDIS_PASSWORD" required:"true"`
	RedisDriver       string   `env:"REDIS_DRIVER" required:"true"`
	RedisMasterName   string   `env:"REDIS_MASTER_NAME"`
	RedisUseDefaultDB int      `env:"REDIS_USE_DEFAULT_DB" required:"true"`
}

func ParseConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
