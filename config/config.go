package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	// Repository is the repository name in OWNER/REPO format
	Repository string `json:"repository"`
	// Snippet is the issue number of the snippet
	Snippet int `json:"snippet"`
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
