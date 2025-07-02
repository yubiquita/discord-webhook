package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/yubiquita/discord-webhook/internal/config"
	"github.com/yubiquita/discord-webhook/internal/webhook"
	"github.com/spf13/cobra"
)

func RunSend(cmd *cobra.Command, configPath, message, webhookURL string, dryRun bool) error {
	if message == "" {
		scanner := bufio.NewScanner(cmd.InOrStdin())
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
		
		if len(lines) == 0 {
			return fmt.Errorf("no message specified")
		}
		
		message = strings.Join(lines, "\n")
	}

	finalWebhookURL := webhookURL
	if finalWebhookURL == "" {
		if configPath == "" {
			configPath = config.GetDefaultConfigPath()
		}

		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load configuration file: %w", err)
		}

		if cfg.WebhookURL == "" {
			return fmt.Errorf("no webhook URL configured. Please specify with --url flag or save in configuration file")
		}

		finalWebhookURL = cfg.WebhookURL
	}

	if dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would send message '%s' to URL '%s'\n", message, finalWebhookURL)
		return nil
	}

	client := webhook.NewClient()
	err := client.SendMessage(finalWebhookURL, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Message sent successfully")
	return nil
}

func RunConfigSet(configPath, key, value string) error {
	if key != "webhook_url" {
		return fmt.Errorf("unsupported configuration key: %s", key)
	}

	if configPath == "" {
		configPath = config.GetDefaultConfigPath()
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %w", err)
	}

	switch key {
	case "webhook_url":
		cfg.WebhookURL = value
	}

	err = cfg.Save(configPath)
	if err != nil {
		return fmt.Errorf("failed to save configuration file: %w", err)
	}

	fmt.Printf("Configuration '%s' set to '%s'\n", key, value)
	return nil
}

func RunConfigGet(configPath, key string) error {
	if configPath == "" {
		configPath = config.GetDefaultConfigPath()
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %w", err)
	}

	if key == "" {
		fmt.Printf("webhook_url: %s\n", cfg.WebhookURL)
		return nil
	}

	switch key {
	case "webhook_url":
		fmt.Println(cfg.WebhookURL)
	default:
		return fmt.Errorf("unsupported configuration key: %s", key)
	}

	return nil
}