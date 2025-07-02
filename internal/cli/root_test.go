package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommand_DisplaysHelpMessageCorrectly(t *testing.T) {
	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("failed to execute help command: %v", err)
	}

	result := output.String()
	if result == "" {
		t.Error("help message is empty")
	}
}

func TestSendCommand_MocksNormalMessageSending(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testWebhookURL := "https://discord.com/api/webhooks/test/token"
	testMessage := "test message"

	config := fmt.Sprintf(`{"webhook_url": "%s"}`, testWebhookURL)
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("failed to create test configuration file: %v", err)
	}

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetArgs([]string{"send", "--config", configPath, "--message", testMessage, "--dry-run"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("failed to execute send command: %v", err)
	}
}

func TestSendCommand_ReturnsErrorWhenMessageNotSpecified(t *testing.T) {
	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetArgs([]string{"send"})

	err := cmd.Execute()
	if err == nil {
		t.Error("no error returned when message is not specified")
	}
}

func TestConfigCommand_WebhookURLConfigurationWorksCorrectly(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testWebhookURL := "https://discord.com/api/webhooks/test/token"

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetArgs([]string{"config", "set", "webhook_url", testWebhookURL, "--config", configPath})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("failed to execute config set command: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("configuration file was not created")
	}
}

func TestSendCommand_ReadsMessageFromStdin(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testWebhookURL := "https://discord.com/api/webhooks/test/token"
	testMessage := "test message from stdin"

	config := fmt.Sprintf(`{"webhook_url": "%s"}`, testWebhookURL)
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("failed to create test configuration file: %v", err)
	}

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetIn(strings.NewReader(testMessage + "\n"))
	cmd.SetArgs([]string{"send", "--config", configPath, "--dry-run"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("failed to execute send command with stdin: %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, testMessage) {
		t.Errorf("output does not contain message from stdin. Output: '%s'", outputStr)
	}
}

func TestSendCommand_FlagMessageTakesPriorityOverStdin(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testWebhookURL := "https://discord.com/api/webhooks/test/token"
	flagMessage := "flag message"
	stdinMessage := "stdin message"

	config := fmt.Sprintf(`{"webhook_url": "%s"}`, testWebhookURL)
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("failed to create test configuration file: %v", err)
	}

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetIn(strings.NewReader(stdinMessage + "\n"))
	cmd.SetArgs([]string{"send", "--config", configPath, "--message", flagMessage, "--dry-run"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("failed to execute send command with flag message priority: %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, flagMessage) {
		t.Errorf("output does not contain flag message. Output: '%s'", outputStr)
	}
	if strings.Contains(outputStr, stdinMessage) {
		t.Errorf("stdin message was used despite flag being specified. Output: '%s'", outputStr)
	}
}

func TestSendCommand_ReturnsErrorWhenStdinIsEmpty(t *testing.T) {
	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"send"})

	err := cmd.Execute()
	if err == nil {
		t.Error("no error returned when stdin is empty")
	}
}