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

type valueSetter map[string]func()

func (vs valueSetter) addSetterFor(name string, setter func()) {
	vs[name] = setter
}

func (vs valueSetter) setValueFor(name string) {
	if setter, exists := vs[name]; exists {
		setter()
	}
}

func addArg[T any](valueSetter *valueSetter, reg argReg[T], cfgValue *configValue[T], name string, value T, usage string) {
	parsedValue := reg(name, value, usage)
	valueSetter.addSetterFor(name, func() {
		cfgValue.Set(*parsedValue, ArgumentSource)
	})
}

func loadFromArgs() *Config {
	var cfg Config

	setter := &valueSetter{}
	addArg[uint64](setter, flag.Uint64, &cfg.Server.Port, "port", 4000, "API server port")

	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		setter.setValueFor(f.Name)
	})

	if !validator.ValidateValueInRange[uint64](cfg.Server.Port.Value, 1024, 65535) {
		log.Fatalf("Invalid port value %d. Should be in-between 1024 and 65535", cfg.Server.Port.Value)
	}
	return &cfg
}
