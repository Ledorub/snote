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

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port configValue[uint64]
}

func New() *Config {
	return loadFromArgs()
}

type argReg[T any] func(name string, value T, usage string) *T

type valueSetters map[string]func()

func (vs valueSetters) addSetterFor(name string, setter func()) {
	vs[name] = setter
}

func (vs valueSetters) setValueFor(name string) {
	if setter, exists := vs[name]; exists {
		setter()
	}
}

func addArg[T any](reg argReg[T], name string, value T, usage string, setters *valueSetters, cfgValue *configValue[T]) {
	parsedValue := reg(name, value, usage)
	setters.addSetterFor(name, func() {
		cfgValue.Set(*parsedValue, ArgumentSource)
	})
}

func loadFromArgs() *Config {
	var cfg Config

	setters := &valueSetters{}
	addArg[uint64](flag.Uint64, "port", 4000, "API server port", setters, &cfg.Server.Port)

	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		setters.setValueFor(f.Name)
	})

	if !validator.ValidateValueInRange[uint64](cfg.Server.Port.Value, 1024, 65535) {
		log.Fatalf("Invalid port value %d. Should be in-between 1024 and 65535", cfg.Server.Port.Value)
	}
	return &cfg
}
