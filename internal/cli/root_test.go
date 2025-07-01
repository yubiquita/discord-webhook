package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func TestSendCommand_標準入力からメッセージを読み取る(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testWebhookURL := "https://discord.com/api/webhooks/test/token"
	testMessage := "標準入力からのテストメッセージ"

	config := fmt.Sprintf(`{"webhook_url": "%s"}`, testWebhookURL)
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("テスト用設定ファイルの作成に失敗しました: %v", err)
	}

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetIn(strings.NewReader(testMessage + "\n"))
	cmd.SetArgs([]string{"send", "--config", configPath, "--dry-run"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("標準入力を使ったsendコマンドの実行に失敗しました: %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, testMessage) {
		t.Errorf("出力に標準入力からのメッセージが含まれていません。出力: '%s'", outputStr)
	}
}

func TestSendCommand_フラグのメッセージが標準入力より優先される(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testWebhookURL := "https://discord.com/api/webhooks/test/token"
	flagMessage := "フラグのメッセージ"
	stdinMessage := "標準入力のメッセージ"

	config := fmt.Sprintf(`{"webhook_url": "%s"}`, testWebhookURL)
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("テスト用設定ファイルの作成に失敗しました: %v", err)
	}

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetIn(strings.NewReader(stdinMessage + "\n"))
	cmd.SetArgs([]string{"send", "--config", configPath, "--message", flagMessage, "--dry-run"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("フラグメッセージ優先のsendコマンドの実行に失敗しました: %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, flagMessage) {
		t.Errorf("出力にフラグのメッセージが含まれていません。出力: '%s'", outputStr)
	}
	if strings.Contains(outputStr, stdinMessage) {
		t.Errorf("フラグが指定されているのに標準入力のメッセージが使用されました。出力: '%s'", outputStr)
	}
}

func TestSendCommand_標準入力が空の場合にエラーを返す(t *testing.T) {
	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"send"})

	err := cmd.Execute()
	if err == nil {
		t.Error("標準入力が空の場合にエラーが返されませんでした")
	}
}