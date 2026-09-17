package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type (
	Config struct {
		Server       ServerConfig `yaml:"server"`
		Services     ServicesURL  `yaml:"services_url"`
		JWTSecretKey string       `yaml:"jwt_secret_key"`
	}

	ServerConfig struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Env  string `yaml:"env"`
	}

	ServicesURL struct {
		User         string `yaml:"user_url"`
		Notification string `yaml:"notification_url"`
	}
)

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	root, err := os.OpenRoot(".")
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}

	f, err := root.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", err)
		}
	}()

	d := yaml.NewDecoder(f)
	if err := d.Decode(config); err != nil {
		return nil, err
	}

	return config, err
}

func getEnv(key, currentVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return currentVal
}
