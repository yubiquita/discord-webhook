package cli

import (
	"fmt"

	"discord-webhook/internal/config"
	"discord-webhook/internal/webhook"
)

func RunSend(configPath, message, webhookURL string, dryRun bool) error {
	if message == "" {
		return fmt.Errorf("メッセージが指定されていません")
	}

	finalWebhookURL := webhookURL
	if finalWebhookURL == "" {
		if configPath == "" {
			configPath = config.GetDefaultConfigPath()
		}

		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
		}

		if cfg.WebhookURL == "" {
			return fmt.Errorf("Webhook URLが設定されていません。--urlフラグで指定するか、設定ファイルに保存してください")
		}

		finalWebhookURL = cfg.WebhookURL
	}

	if dryRun {
		fmt.Printf("Dry run: メッセージ「%s」をURL「%s」に送信する予定です\n", message, finalWebhookURL)
		return nil
	}

	client := webhook.NewClient()
	err := client.SendMessage(finalWebhookURL, message)
	if err != nil {
		return fmt.Errorf("メッセージの送信に失敗しました: %w", err)
	}

	fmt.Println("メッセージを正常に送信しました")
	return nil
}

func RunConfigSet(configPath, key, value string) error {
	if key != "webhook_url" {
		return fmt.Errorf("サポートされていない設定キー: %s", key)
	}

	if configPath == "" {
		configPath = config.GetDefaultConfigPath()
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	switch key {
	case "webhook_url":
		cfg.WebhookURL = value
	}

	err = cfg.Save(configPath)
	if err != nil {
		return fmt.Errorf("設定ファイルの保存に失敗しました: %w", err)
	}

	fmt.Printf("設定「%s」を「%s」に設定しました\n", key, value)
	return nil
}

func RunConfigGet(configPath, key string) error {
	if configPath == "" {
		configPath = config.GetDefaultConfigPath()
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	if key == "" {
		fmt.Printf("webhook_url: %s\n", cfg.WebhookURL)
		return nil
	}

	switch key {
	case "webhook_url":
		fmt.Println(cfg.WebhookURL)
	default:
		return fmt.Errorf("サポートされていない設定キー: %s", key)
	}

	return nil
}