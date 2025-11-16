package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DbConfig
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" env-default:"8080"`
}

type DbConfig struct {
	Driver   string `env:"DB_DRIVER" env-required:"true"`
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	DB       string `env:"DB_NAME" env-required:"true"`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found or error loading it:", err)
	}

	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("read env: %v", err)
	}

	return &cfg
}
