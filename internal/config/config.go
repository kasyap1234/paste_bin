package config

import "os"

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	UseSSl   bool
}

type LoggerConfig struct {
	Level string
}

func LoadLoggerConfig() *LoggerConfig {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	return &LoggerConfig{Level: level}
}
