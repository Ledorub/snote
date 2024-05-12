package config

import "flag"

type args struct {
	port       uint64
	configFile string
}

func (l *Loader) loadArgs() *args {
	return loadArgs()
}

func loadArgs() *args {
	a := args{}
	flag.Uint64Var(&a.port, "port", 4000, "API server port")
	flag.StringVar(&a.configFile, "config-file", "", "Path to a config file")
	flag.Parse()

	return &a
}
