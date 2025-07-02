package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "discord-webhook",
		Short: "CLI tool for sending messages via Discord Webhook",
		Long: `A simple CLI tool for sending messages via Discord Webhook.
Manage Webhook URLs in configuration files and send messages from the command line.`,
	}

	rootCmd.AddCommand(NewSendCommand())
	rootCmd.AddCommand(NewConfigCommand())

	return rootCmd
}

func NewSendCommand() *cobra.Command {
	var (
		configPath string
		message    string
		webhookURL string
		dryRun     bool
	)

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send message to Discord Webhook",
		Long:  "Send the specified message to Discord Webhook.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSend(cmd, configPath, message, webhookURL, dryRun)
		},
	}

	sendCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	sendCmd.Flags().StringVarP(&message, "message", "m", "", "Message to send (reads from stdin if not specified)")
	sendCmd.Flags().StringVarP(&webhookURL, "url", "u", "", "Webhook URL (takes priority over configuration file)")
	sendCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Test execution without actually sending")

	return sendCmd
}

func NewConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long: `Manage configurations such as webhook URLs.

Available configuration keys:
  webhook_url    Discord Webhook URL
                 Format: https://discord.com/api/webhooks/{id}/{token}
                 Usage: Discord Webhook URL for message destination

Configuration file:
  Default: ~/.discord-webhook/config.json
  Custom: Can be specified with --config flag

Usage examples:
  discord-webhook config set webhook_url https://discord.com/api/webhooks/...
  discord-webhook config get webhook_url
  discord-webhook config get`,
	}

	configCmd.AddCommand(NewConfigSetCommand())
	configCmd.AddCommand(NewConfigGetCommand())

	return configCmd
}

func NewConfigSetCommand() *cobra.Command {
	var configPath string

	setCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set configuration value",
		Long:  "Set configuration value. Currently only webhook_url is supported.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSet(configPath, args[0], args[1])
		},
	}

	setCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")

	return setCmd
}

func NewConfigGetCommand() *cobra.Command {
	var configPath string

	getCmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Get configuration value",
		Long:  "Display configuration values. Shows all configurations if key is not specified.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := ""
			if len(args) > 0 {
				key = args[0]
			}
			return runConfigGet(configPath, key)
		},
	}

	getCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")

	return getCmd
}

func runSend(cmd *cobra.Command, configPath, message, webhookURL string, dryRun bool) error {
	return RunSend(cmd, configPath, message, webhookURL, dryRun)
}

func runConfigSet(configPath, key, value string) error {
	return RunConfigSet(configPath, key, value)
}

func runConfigGet(configPath, key string) error {
	return RunConfigGet(configPath, key)
}