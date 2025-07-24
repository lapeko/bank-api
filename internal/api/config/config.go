package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type config struct {
	AppPort      int    `env:"APP_PORT" envDefault:"3000"`
	PostgresUrl  string `env:"POSTGRES_URL,required"`
	JwtSecretKey string `env:"JWT_SECRET_KEY,required"`
}

var cfg *config

func Get() config {
	if cfg != nil {
		return *cfg
	}

	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Config parse error: %v", err)
	}
	return *cfg
}
