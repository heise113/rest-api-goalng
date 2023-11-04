package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"development"`
	Storage    `yaml:"storage"`
	HTTPServer `yaml:"http_server"`
}

type Storage struct {
	AddressDB string `yaml:"address_db"`
	Login     string `yaml:"login"`
	Pass      string `yaml:"pass"`
	NameDB    string `yaml:"name_db"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	// (Я) Создаю переменную окружения
	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	// Получаем путь до конфиг-файла из env-переменной CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// Проверяем существование конфиг-файла
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	// Читаем конфиг-файл и заполняем нашу структуру
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	fmt.Println("cfg: ", cfg)
	return &cfg
}
