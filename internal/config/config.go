package config

import (
	"flag"
	"github.com/ledorub/snote-api/internal/validator"
	"log"
)

type sourceType string

const (
	SourceUnset    sourceType = ""
	FileSource     sourceType = "file"
	ArgumentSource sourceType = "argument"
)

type Config struct {
	Server ServerConfig
}

func (cfg *Config) load() {

	parsedArgs := loadArgs()

	if !validator.ValidateValueInRange[uint64](parsedArgs.port, 1024, 65535) {
		log.Fatalf("Invalid port value %d. Should be in-between 1024 and 65535", parsedArgs.port)
	}

	setters := configValueSetters{}
	mapArgsToConfigValues(setters, parsedArgs, cfg)
}

type ServerConfig struct {
	Port configValue[uint64]
}

func Load() *Config {
	cfg := &Config{}
	cfg.load()
	return cfg
}

type configValue[T any] struct {
	Value  T
	Source sourceType
}

func (cv *configValue[T]) Set(value T, src sourceType) {
	if cv.Source == SourceUnset {
		cv.Value = value
		cv.Source = src
	}
}

type configValueSetters map[string]func()

func (s configValueSetters) addSetterFor(name string, setter func()) {
	s[name] = setter
}

func (s configValueSetters) setValueFor(name string) {
	if setter, exists := s[name]; exists {
		setter()
	}
}

func (s configValueSetters) setValueForAll() {
	for _, setter := range s {
		setter()
	}
}

type args struct {
	port uint64
}

func loadArgs() *args {
	a := args{}
	flag.Uint64Var(&a.port, "port", 4000, "API server port")
	flag.Parse()

	return &a
}

func mapToConfigValue[T any](mp configValueSetters, name string, src sourceType, from *T, to *configValue[T]) {
	mp.addSetterFor(name, func() {
		to.Set(*from, src)
	})
}

func mapArgsToConfigValues(mp configValueSetters, a *args, cfg *Config) {
	mapToConfigValue[uint64](mp, "port", ArgumentSource, &a.port, &cfg.Server.Port)
}
