package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_保存と読み込みが正常に動作する(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	config := &Config{
		WebhookURL: "https://discord.com/api/webhooks/test/token",
	}

	err := config.Save(configPath)
	if err != nil {
		t.Fatalf("設定の保存に失敗しました: %v", err)
	}

	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	if loadedConfig.WebhookURL != config.WebhookURL {
		t.Errorf("期待されるWebhook URL: %s, 実際: %s", config.WebhookURL, loadedConfig.WebhookURL)
	}
}

func TestConfig_存在しないファイルからの読み込みでデフォルト設定を返す(t *testing.T) {
	nonExistentPath := "/tmp/nonexistent_config.json"

	config, err := Load(nonExistentPath)
	if err != nil {
		t.Fatalf("存在しないファイルからの読み込みでエラーが発生しました: %v", err)
	}

	if config.WebhookURL != "" {
		t.Errorf("デフォルト設定でWebhook URLが空文字列でありません: %s", config.WebhookURL)
	}
}

func TestConfig_不正なJSONファイルからの読み込みでエラーを返す(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid_config.json")

	err := os.WriteFile(configPath, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("テスト用不正JSONファイルの作成に失敗しました: %v", err)
	}

	_, err = Load(configPath)
	if err == nil {
		t.Error("不正なJSONファイルに対してエラーが返されませんでした")
	}
}

func TestConfig_デフォルト設定パスが正しく取得される(t *testing.T) {
	path := GetDefaultConfigPath()

	if path == "" {
		t.Error("デフォルト設定パスが空文字列です")
	}

	if !filepath.IsAbs(path) {
		t.Error("デフォルト設定パスが絶対パスではありません")
	}
}