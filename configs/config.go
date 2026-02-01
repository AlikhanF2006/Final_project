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
}
