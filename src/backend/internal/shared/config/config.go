package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Kafka    KafkaConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type KafkaConfig struct {
	Brokers      []string
	Topic        string
	WorkerCount  int
	PollInterval time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() (*Config, error) {
	workerCount, err := getEnvInt("KAFKA_WORKER_COUNT", 3)
	if err != nil {
		return nil, fmt.Errorf("KAFKA_WORKER_COUNT: %w", err)
	}

	pollInterval, err := getEnvDuration("KAFKA_POLL_INTERVAL", 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("KAFKA_POLL_INTERVAL: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Kafka: KafkaConfig{
			Brokers:      strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			Topic:        getEnv("KAFKA_TOPIC", "email"),
			WorkerCount:  workerCount,
			PollInterval: pollInterval,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "personal_page"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("must be an integer, got %q", v)
	}
	return n, nil
}

func getEnvDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("must be a duration (e.g. 2s, 500ms), got %q", v)
	}
	return d, nil
}
