package config

import (
	"bufio"
	"fmt"
	"os"
)

type configFile struct {
	Server configFileServer
	DB     configFileDB
}

type configFileServer struct {
	Port uint64 `yaml:"port"`
}

type configFileDB struct {
	Host     string       `yaml:"host"`
	Port     uint64       `yaml:"port"`
	Name     string       `yaml:"name"`
	User     string       `yaml:"user"`
	Password secretString `yaml:"password"`
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

func getFileReader(path string) (*bufio.Reader, error) {
	f, err := openFile(path)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(f), nil
}

func openFile(path string) (*os.File, error) {
	return os.Open(path)
}
