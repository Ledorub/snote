package config

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ledorub/snote-api/internal/encdec"
	"github.com/ledorub/snote-api/internal/validator"
	"io"
	"os"
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

func (cfg *Config) checkErrors() error {
	if !validator.ValidateValueInRange[uint64](cfg.Server.Port.Value, 1024, 65535) {
		return fmt.Errorf("invalid port value %d. Should be in-between 1024 and 65535", cfg.Server.Port.Value)
	}
	return nil
}

type ServerConfig struct {
	Port configValue[uint64]
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
	port       uint64
	configFile string
}

func loadArgs() *args {
	a := args{}
	flag.Uint64Var(&a.port, "port", 4000, "API server port")
	flag.StringVar(&a.configFile, "config-file", "", "Path to a config file")
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

func mapConfigFileToConfigValues(mp configValueSetters, cfgF *configFile, cfg *Config) {
	mapToConfigValue[uint64](mp, "port", FileSource, &cfgF.Server.Port, &cfg.Server.Port)
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

	if l.shouldLoadArgs {
		loadedArgs := l.loadArgs()
		mapArgsToConfigValues(setters, loadedArgs, cfg)

		if l.configFile == "" {
			l.configFile = loadedArgs.configFile
		}
	}

	if l.configFile != "" {
		fileCfg, err := l.loadFile()
		if err != nil {
			return nil, err
		}
		mapConfigFileToConfigValues(setters, fileCfg, cfg)
	}

	setters.setValueForAll()
	if err := cfg.checkErrors(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (l *Loader) loadArgs() *args {
	return loadArgs()
}

func (l *Loader) loadFile() (*configFile, error) {
	reader, err := getFileReader(l.configFile)
	if err != nil {
		return nil, fmt.Errorf("config file loader: %v", err)
	}

	fileConfig := &configFile{}
	err = l.configFileDecoder.Decode(reader, fileConfig)
	if err != nil {
		return nil, fmt.Errorf("config file loader: %v", err)
	}
	return fileConfig, nil
}

func openFile(path string) (*os.File, error) {
	return os.Open(path)
}

func getFileReader(path string) (*bufio.Reader, error) {
	f, err := openFile(path)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(f), nil
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

type configFileServer struct {
	Port uint64
}

type configFile struct {
	Server configFileServer
}
