package config

import "os"


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
