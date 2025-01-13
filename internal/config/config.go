package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const cfg_filename string = ".gatorconfig.json"

func getCfgFilepath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filepath := home + "/" + cfg_filename
	return filepath, nil
}

func Read() (Config, error) {
	filepath, err := getCfgFilepath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUser = user

	filepath, err := getCfgFilepath()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, bytes, 0600); err != nil {
		return err
	}

	return nil
}
