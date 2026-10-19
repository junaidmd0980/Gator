package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := GetConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading file: %w", err)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error decoding data: %w", err)
	}
	return config, nil
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error encoding data: %w", err)
	}
	path, err := GetConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("Error writing file: %w", err)
	}
	return nil
}

func (cfg Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(cfg)
}

const configFileName = "/.gatorconfig.json"

func GetConfigFilePath() (string, error) {
	home_path, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting home path: %w", err)
	}
	path := home_path + configFileName
	return path, nil
}
