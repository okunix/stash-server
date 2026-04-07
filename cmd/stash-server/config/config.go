package config

import (
	"errors"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslMode"`
}

type Config struct {
	Addr           string         `yaml:"addr"`
	LogFile        string         `yaml:"logFile"`
	PostgresConfig PostgresConfig `yaml:"postgres"`
}

func DefaultConfig() Config {
	return Config{
		Addr:    "0.0.0.0:7878",
		LogFile: "",
		PostgresConfig: PostgresConfig{
			Host:     "localhost",
			Port:     "5432",
			Database: "stash_db",
			User:     "postgres",
			Password: "postgres",
			SSLMode:  "disable",
		},
	}
}

func ReadFromFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return DefaultConfig(), err
	}
	defer file.Close()
	return Read(file)
}

func Read(r io.Reader) (Config, error) {
	conf = DefaultConfig()
	err := yaml.NewDecoder(r).Decode(&conf)
	if err != nil && !errors.Is(err, io.EOF) {
		return conf, err
	}
	return conf, nil
}

var conf Config

func Conf() Config {
	return conf
}
