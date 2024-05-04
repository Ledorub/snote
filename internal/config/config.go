package config

import (
	"flag"
	"github.com/ledorub/snote-api/internal/validator"
	"log"
)

type source string

const (
	FileSource     source = "file"
	ArgumentSource source = "argument"
)

type configValue[T any] struct {
	Value  T
	Source source
}

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port configValue[uint64]
}

func New() *Config {
	return loadFromArgs()
}

func loadFromArgs() *Config {
	var cfg Config

	flag.Uint64Var(&cfg.Server.Port.Value, "port", 4000, "API server port")
	cfg.Server.Port.Source = FileSource
	flag.Parse()

	if !validator.ValidateValueInRange[uint64](cfg.Server.Port.Value, 1024, 65535) {
		log.Fatalf("Invalid port value %d. Should be in-between 1024 and 65535", cfg.Server.Port)
	}
	return &cfg
}
