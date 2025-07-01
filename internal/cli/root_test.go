package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRootCommand_正常にヘルプメッセージを表示する(t *testing.T) {
	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("ヘルプコマンドの実行に失敗しました: %v", err)
	}

	result := output.String()
	if result == "" {
		t.Error("ヘルプメッセージが空です")
	}
}

func TestSendCommand_正常なメッセージ送信をモックする(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testWebhookURL := "https://discord.com/api/webhooks/test/token"
	testMessage := "テストメッセージ"

	config := fmt.Sprintf(`{"webhook_url": "%s"}`, testWebhookURL)
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("テスト用設定ファイルの作成に失敗しました: %v", err)
	}

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetArgs([]string{"send", "--config", configPath, "--message", testMessage, "--dry-run"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("sendコマンドの実行に失敗しました: %v", err)
	}
}

func TestSendCommand_メッセージが指定されていない場合にエラーを返す(t *testing.T) {
	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetArgs([]string{"send"})

	err := cmd.Execute()
	if err == nil {
		t.Error("メッセージが指定されていない場合にエラーが返されませんでした")
	}
}

func TestConfigCommand_webhook_url設定が正常に動作する(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testWebhookURL := "https://discord.com/api/webhooks/test/token"

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetArgs([]string{"config", "set", "webhook_url", testWebhookURL, "--config", configPath})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("config setコマンドの実行に失敗しました: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("設定ファイルが作成されませんでした")
	}
}