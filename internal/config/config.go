package config

import (
	"fmt"
	"github.com/ledorub/snote-api/internal/encdec"
	"github.com/ledorub/snote-api/internal/validator"
	"io"
)

type sourceType string

const (
	SourceUnset    sourceType = ""
	FileSource     sourceType = "file"
	ArgumentSource sourceType = "argument"
)

type Config struct {
	Source ConfigSource `yaml:"source"`
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

func (cfg *Config) checkErrors() error {
	if !validator.ValidateValueInRange[uint64](cfg.Server.Port.Value, 1024, 65535) {
		return fmt.Errorf("invalid port value %d. Should be in-between 1024 and 65535", cfg.Server.Port.Value)
	}
	if !validator.ValidateValueInRange[uint64](cfg.DB.Port.Value, 1024, 65535) {
		return fmt.Errorf("invalid DB port value %d. Should be in-between 1024 and 65535", cfg.DB.Port)
	}
	return nil
}

func (cfg *Config) Pretty() (string, error) {
	enc := encdec.NewYAMLEncoder()
	encoded, err := enc.Encode(cfg)
	if err != nil {
		return "", fmt.Errorf("config prettier: %w", err)
	}
	return string(encoded), nil
}

type ConfigSource struct {
	ParseArgs bool   `yaml:"parseArgs"`
	File      string `yaml:"file"`
}

type ServerConfig struct {
	Port configValue[uint64] `yaml:"port"`
}

type DBConfig struct {
	Host     configValue[string]       `yaml:"host"`
	Port     configValue[uint64]       `yaml:"port"`
	Name     configValue[string]       `yaml:"name"`
	User     configValue[string]       `yaml:"user"`
	Password configValue[secretString] `yaml:"password"`
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

type Loader struct {
	shouldLoadArgs    bool
	configFile        string
	configFileDecoder configFileDecoder
}

func NewLoader(opts ...LoaderOpt) *Loader {
	loader := &Loader{}
	for _, opt := range opts {
		opt(loader)
	}

	if loader.configFileDecoder == nil {
		loader.configFileDecoder = encdec.NewYAMLDecoder()
	}
	return loader
}

func (l *Loader) Load() (*Config, error) {
	cfg := &Config{}
	setters := configValueSetters{}
	mapper := newValueMapper(setters, cfg)

	if l.shouldLoadArgs {
		loadedArgs := l.loadArgs()
		cfg.Source.ParseArgs = true
		mapper.mapArgsToConfigValues(loadedArgs)

		if l.configFile == "" {
			l.configFile = loadedArgs.configFile
		}
	}

	if l.configFile != "" {
		fileCfg, err := l.loadFile()
		if err != nil {
			return nil, err
		}
		cfg.Source.File = l.configFile
		mapper.mapConfigFileToConfigValues(fileCfg)
	}

	setters.setValueForAll()
	if err := cfg.checkErrors(); err != nil {
		return nil, err
	}
	return cfg, nil
}

type LoaderOpt func(l *Loader)

func LoadArgs() LoaderOpt {
	return func(l *Loader) {
		l.shouldLoadArgs = true
	}
}

func LoadFile(path string, decoder configFileDecoder) LoaderOpt {
	return func(l *Loader) {
		l.configFile = path
		l.configFileDecoder = decoder
	}
}

type configFileDecoder interface {
	Decode(data io.Reader, dst any) error
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

type valueMapper struct {
	setters configValueSetters
	config  *Config
}

func newValueMapper(setters configValueSetters, config *Config) *valueMapper {
	return &valueMapper{setters: setters, config: config}
}

func (m *valueMapper) mapArgsToConfigValues(a *args) {
	src := ArgumentSource
	mapToConfigValue[uint64](m.setters, "port", src, &a.port, &m.config.Server.Port)
}

func (m *valueMapper) mapConfigFileToConfigValues(cfgF *configFile) {
	src := FileSource
	mapToConfigValue[uint64](m.setters, "port", src, &cfgF.Server.Port, &m.config.Server.Port)
	mapToConfigValue[string](m.setters, "db_host", src, &cfgF.DB.Host, &m.config.DB.Host)
	mapToConfigValue[uint64](m.setters, "db_port", src, &cfgF.DB.Port, &m.config.DB.Port)
	mapToConfigValue[string](m.setters, "db_name", src, &cfgF.DB.Name, &m.config.DB.Name)
	mapToConfigValue[string](m.setters, "db_user", src, &cfgF.DB.User, &m.config.DB.User)
	mapToConfigValue[secretString](m.setters, "db_password", src, &cfgF.DB.Password, &m.config.DB.Password)
}

func mapToConfigValue[T any](mp configValueSetters, name string, src sourceType, from *T, to *configValue[T]) {
	mp.addSetterFor(name, func() {
		to.Set(*from, src)
	})
}
