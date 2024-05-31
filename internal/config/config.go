package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HttpServer `yaml:"http_server"`
	Storage    `yaml:"storage"`
	NasaApi    `yaml:"nasa_api_models"`
	CtxTimeout time.Duration `yaml:"ctx_timeout" env-default:"5s"`
}

type HttpServer struct {
	Host         string        `yaml:"host" env-default:"localhost:8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Storage struct {
	DbUri string `yaml:"db_uri" env-required:"true"`
}

type NasaApi struct {
	Host  string `yaml:"host" env-required:"true"`
	Token string `yaml:"token" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check config if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return &cfg
}
