package config

import (
	"fmt"
	"io/ioutil"
	"net/smtp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	} `yaml:"db"`

	Redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}

	SMTP struct {
		Client   *smtp.Client
		Host     string `yaml:"SMTP_host"`
		Port     string `yaml:"SMTP_port"`
		Email    string `yaml:"SMTP_email"`
		Password string `yaml:"SMTP_password"`
	}
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	return &config, nil
}
