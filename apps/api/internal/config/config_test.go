package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	config, _ := LoadConfig()
	if config.DatabaseURL == "" {
		t.Error("Expected DatabaseURL to be set")
	}
}
