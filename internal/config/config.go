package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name     string       `yaml:"Name"`
	Host     string       `yaml:"Host"`
	Port     int          `yaml:"Port"`
	DBPath   string       `yaml:"DBPath"`
	Auth115  Auth115Config `yaml:"Auth115"`
	Log      LogConfig     `yaml:"Log"`
}

type Auth115Config struct {
	DownloadPath string `yaml:"DownloadPath"`
	AccessToken  string `yaml:"AccessToken"`
	RefreshToken string `yaml:"RefreshToken"`
}

type LogConfig struct {
	Level  string `yaml:"Level"`
	Format string `yaml:"Format"`
}

func DefaultConfig() *Config {
	return &Config{
		Name: "115quick",
		Host: "0.0.0.0",
		Port: 8889,
		DBPath: "data/115quick.db",
		Auth115: Auth115Config{
			DownloadPath: "data",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
