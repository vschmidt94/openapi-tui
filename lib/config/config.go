package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vschmidt94/openapi-tui/types"
	"os"
	"path/filepath"
	"sort"
)

type Config struct {
	Sites []types.Site `mapstructure:"sites"`
}

func LoadConfig() (*Config, error) {
	var cfg = new(Config)
	cwd, err := os.Getwd()
	fmt.Println(cwd)
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	configPath := filepath.Join(cwd, "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		viper.AddConfigPath(filepath.Dir(configPath))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Parse the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	// Sort the instances by name
	sort.Slice(cfg.Sites, func(i, j int) bool {
		return cfg.Sites[i].Name < cfg.Sites[j].Name
	})

	return cfg, nil
}
