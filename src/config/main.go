package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"dev"`
	GitHubOauth
}

func LoadEnv() *Config {
	config := Config{}
	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatal("Failed to LoadEnv:", err)
	}

	if config.Env == "dev" {
		if err := cleanenv.ReadConfig(".env", &config); err != nil {
			log.Fatal("Failed to load .env: ", err)
		}
	}

	return &config
}
