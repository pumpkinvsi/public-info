package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Kafka KafkaConfig
	Smtp  SmtpConfig
	Loki    LokiConfig
	Metrics MetricsConfig
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type SmtpConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type LokiConfig struct {
	URL string
	AppLabel string
}

type MetricsConfig struct {
	Port int
}


func Load() (*Config, error) {
	smtpPort, err := getEnvInt("SMTP_PORT", 587)
	if err != nil {
		return nil, fmt.Errorf("SMTP_PORT: %w", err)
	}
	metricsPort, err := getEnvInt("METRICS_PORT", 9090)
	if err != nil {
		return nil, fmt.Errorf("METRICS_PORT: %w", err)
	}

	cfg := &Config{
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			Topic:   getEnv("KAFKA_TOPIC", "email"),
			GroupID: getEnv("KAFKA_GROUP_ID", "email-sender"),
		},
		Smtp: SmtpConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     smtpPort,
			User:     getEnv("SMTP_USER", "user"),
			Password: getEnv("SMTP_PASSWORD", "password"),
		},
		Loki: LokiConfig{
			URL: getEnv("LOKI_URL", "http://loki:3100"),
			AppLabel: getEnv("LOKI_APP_LABEL", "email-sender"),
		},
		Metrics: MetricsConfig{
			Port: metricsPort,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if len(c.Kafka.Brokers) == 0 {
		return fmt.Errorf("KAFKA_BROKERS environment variable is required")
	}
	if strings.TrimSpace(c.Kafka.Topic) == "" {
		return fmt.Errorf("KAFKA_TOPIC environment variable is required")
	}
	if strings.TrimSpace(c.Kafka.GroupID) == "" {
		return fmt.Errorf("KAFKA_GROUP_ID environment variable is required")
	}
	if strings.TrimSpace(c.Smtp.Host) == "" {
		return fmt.Errorf("SMTP_HOST environment variable is required")
	}
	if c.Smtp.Port == 0 {
		return fmt.Errorf("SMTP_PORT environment variable is required")
	}
	if strings.TrimSpace(c.Smtp.User) == "" {
		return fmt.Errorf("SMTP_USER environment variable is required")
	}
	if strings.TrimSpace(c.Smtp.Password) == "" {
		return fmt.Errorf("SMTP_PASSWORD environment variable is required")
	}
	if strings.TrimSpace(c.Loki.URL) == "" {
		return fmt.Errorf("LOKI_URL environment variable is required")
	}
	if strings.TrimSpace(c.Loki.AppLabel) == "" {
		return fmt.Errorf("LOKI_APP_LABEL environment variable is required")
	}
	if c.Metrics.Port == 0 {
		return fmt.Errorf("METRICS_PORT environment variable is required")
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