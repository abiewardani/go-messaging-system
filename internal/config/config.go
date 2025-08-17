package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DatabaseURL string `json:"database_url"`
	RabbitMQURL string `json:"rabbitmq_url"`
	Concurrency  int    `json:"concurrency"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}