package configs

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`

	TMDB struct {
		ApiKey string `yaml:"api_key"`
	} `yaml:"tmdb"`
}

var AppConfig Config

func LoadConfig() {
	file, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatal("cannot read config.yaml:", err)
	}

	if err := yaml.Unmarshal(file, &AppConfig); err != nil {
		log.Fatal("cannot parse config.yaml:", err)
	}

	if AppConfig.Database.URL == "" {
		log.Fatal("database.url is empty")
	}
	if AppConfig.TMDB.ApiKey == "" {
		log.Fatal("tmdb.api_key is empty")
	}
}
