package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	User     string        `yaml:"user" env-default:""`
	Pass     string        `yaml:"pass" env-default:""`
	Host     string        `yaml:"host" env-default:"localhost"`
	HTTPAddr string        `yaml:"http_addr" env-default:":8080"`
	Port     int           `yaml:"port" env-default:"5432"`
	Name     string        `yaml:"name" env-default:""`
	MaxConns int32         `yaml:"max_conns" env-default:"5"`
	MinConns int32         `yaml:"min_conns" env-default:"2"`
	Timeout  time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
	// читаем путь к конфигу из переменной окружения, при пустом значении используем дефолт
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Println("CONFIG_PATH is empty, using default local config")
		configPath = "config/local.yaml"
	}

	// проверяем существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config: file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("config: cannot read config %s: %v", configPath, err)
	}

	return &cfg
}
