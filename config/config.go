package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Nats struct {
		Host string `yaml:"host"`
	}
	JWT struct {
		SecretKey string `yaml:"secretKey"`
	}
}

func LoadConfig() (*Config, error) {
	stage := getStage()
	configFile := fmt.Sprintf("config-%s.yml", stage)

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}

	configPath := filepath.Join(cwd, "config", configFile)

	yamlData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	config := &Config{}

	err = yaml.Unmarshal(yamlData, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return config, nil
}

func getStage() string {
	stage, exists := os.LookupEnv("STAGE")
	if !exists {
		stage = "dev"
	}

	return stage
}
