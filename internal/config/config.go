package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDirName  = ".ego"
	configFileName = "config.json"
)

type Config struct {
	Active  string `json:"active"`            // name of active vault
	RootDir string `json:"rootDir,omitempty"` // base directory for vaults
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get home directory: %w", err)
	}
	cfgPath := filepath.Join(home, configDirName, configFileName)
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get home directory: %w", err)
	}
	cfgDir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(cfgDir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	cfgPath := filepath.Join(cfgDir, configFileName)
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(cfgPath, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
