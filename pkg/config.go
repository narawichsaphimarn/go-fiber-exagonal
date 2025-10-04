package pkg

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() (*Config, error) {
	var config Config
	configFile, err := os.ReadFile("configs/app.yaml")
	if err != nil {
		return &config, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}

type Config struct {
	App  AppConfig  `yaml:"app"`
	DB   DBConfig   `yaml:"db"`
	Auth AuthConfig `yaml:"auth"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Options  string `yaml:"options"`
}

type AuthConfig struct {
	JWT JWTConfig `yaml:"jwt"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}
