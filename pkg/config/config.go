package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	// Owner is the owner name of the repository
	Owner string `json:"owner"`
	// Repository is the repository name
	Repository string `json:"repository"`
	// Utilities is the issue number of utilities
	Utilities int `json:"utilities"`
	// DataDir is the directory to store data
	DataDir string `json:"datadir"`
}

func LoadConfig(filepath string) (*Config, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = json.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
