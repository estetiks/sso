package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env"`
	StoragePath string        `yaml:"storage_path"`
	TokenTTL    time.Duration `yaml:"token_ttl"`
	GRPC        GRPCConfig    `yaml:"gRPC"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	ConfigPath := fetchConfigPath()

	if ConfigPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		panic("config file does not exist " + ConfigPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(ConfigPath, &cfg); err != nil {
		panic("config path is empty " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
