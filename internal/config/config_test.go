package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_SaveAndLoadWorksCorrectly(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	config := &Config{
		WebhookURL: "https://discord.com/api/webhooks/test/token",
	}

	err := config.Save(configPath)
	if err != nil {
		t.Fatalf("failed to save configuration: %v", err)
	}

	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("failed to load configuration: %v", err)
	}

	if loadedConfig.WebhookURL != config.WebhookURL {
		t.Errorf("expected Webhook URL: %s, actual: %s", config.WebhookURL, loadedConfig.WebhookURL)
	}
}

func TestConfig_LoadFromNonExistentFileReturnsDefaultConfig(t *testing.T) {
	nonExistentPath := "/tmp/nonexistent_config.json"

	config, err := Load(nonExistentPath)
	if err != nil {
		t.Fatalf("error occurred when loading from non-existent file: %v", err)
	}

	if config.WebhookURL != "" {
		t.Errorf("default configuration Webhook URL is not empty string: %s", config.WebhookURL)
	}
}

func TestConfig_LoadFromInvalidJSONFileReturnsError(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid_config.json")

	err := os.WriteFile(configPath, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("failed to create test invalid JSON file: %v", err)
	}

	_, err = Load(configPath)
	if err == nil {
		t.Error("no error returned for invalid JSON file")
	}
}

func TestConfig_DefaultConfigPathIsRetrievedCorrectly(t *testing.T) {
	path := GetDefaultConfigPath()

	if path == "" {
		t.Error("default configuration path is empty string")
	}

	if !filepath.IsAbs(path) {
		t.Error("default configuration path is not absolute path")
	}
}