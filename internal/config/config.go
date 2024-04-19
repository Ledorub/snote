package config

import (
	"flag"
	"github.com/ledorub/snote-api/internal/validator"
	"log"
)

type Config struct {
	ServerConfig
}

type ServerConfig struct {
	Port int
}

func New() *Config {
	return loadFromArgs()
}

func loadFromArgs() *Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.Parse()

	if !validator.ValidateValueInRange[int](cfg.Port, 1024, 65535) {
		log.Fatalf("Invalid port value %d. Should be in-between 1024 and 65535", cfg.Port)
	}
	return &cfg
}
