package config

import (
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

type SQLiteConfig struct {
	DbPath string `yaml:"dbPath"`
}

type Config struct {
	Addr         string       `yaml:"addr"`
	LogFile      string       `yaml:"logFile"`
	SQLiteConfig SQLiteConfig `yaml:"sqlite"`
}

func DefaultConfig() Config {
	return Config{
		Addr:    "0.0.0.0:80",
		LogFile: "",
		SQLiteConfig: SQLiteConfig{
			DbPath: "stash.db",
		},
	}
}

func ReadFromFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return DefaultConfig(), err
	}
	return Read(file)
}

func Read(r io.Reader) (Config, error) {
	conf = DefaultConfig()
	err := yaml.NewDecoder(r).Decode(&conf)
	return conf, err
}

var conf Config

func Conf() Config {
	return conf
}
