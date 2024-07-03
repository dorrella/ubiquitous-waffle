package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

const config_path string = "/config/config.yaml"

type Database struct {
	Password string `yaml: "password"`
	User     string `yaml: "user"`
	Database string `yaml: "database"`
	Hostname string `yaml: "hostname"`
	Port     int    `yaml: "port" `
}

type Collector struct {
	Url  string `yaml: "url"`
	Port int    `yaml: "port"`
}
type Telemetry struct {
	Enabled   bool      `yaml: "enabled"`
	Tracing   bool      `yaml: "tracing"`
	Metrics   bool      `yaml: "metrics"`
	Logs      bool      `yaml: "logs"`
	Collector Collector `yaml: "collector"`
}

type Service struct {
	Name string `yaml: "name"`
}

type Config struct {
	Service   Service   `yaml: "database"`
	Database  Database  `yaml: "database"`
	Telemetry Telemetry `yaml: "telemetry"`
}

func LoadConfig() *Config {
	config := &Config{}
	yaml_file, err := os.ReadFile(config_path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yaml_file, config)
	if err != nil {
		panic(err)
	}
	return config
}

func TestConfig() *Config {
	return &Config{
		Service: Service{
			Name: "test-service",
		},
		Database: Database{
			Password: os.Getenv("POSTGRES_PASSWORD"),
			User:     os.Getenv("POSTGRES_USER"),
			Database: os.Getenv("POSTGRES_DB"),
			Hostname: "postgresql",
			Port:     5432,
		},
		Telemetry: Telemetry{
			Enabled:   false,
			Tracing:   false,
			Metrics:   false,
			Logs:      false,
			Collector: Collector{},
		},
	}
}

func (c *Config) GetPostgresUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Hostname,
		c.Database.Port,
		c.Database.Database)
}
