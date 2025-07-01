package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "discord-webhook",
		Short: "Discord Webhook経由でメッセージを送信するCLIツール",
		Long: `Discord Webhook経由でメッセージを送信するシンプルなCLIツールです。
設定ファイルでWebhook URLを管理し、コマンドラインからメッセージを送信できます。`,
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
		Short: "Discord Webhookにメッセージを送信",
		Long:  "指定されたメッセージをDiscord Webhookに送信します。",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSend(cmd, configPath, message, webhookURL, dryRun)
		},
	}

	sendCmd.Flags().StringVarP(&configPath, "config", "c", "", "設定ファイルのパス")
	sendCmd.Flags().StringVarP(&message, "message", "m", "", "送信するメッセージ（未指定の場合は標準入力から読み取り）")
	sendCmd.Flags().StringVarP(&webhookURL, "url", "u", "", "Webhook URL（設定ファイルより優先されます）")
	sendCmd.Flags().BoolVar(&dryRun, "dry-run", false, "実際に送信せずにテスト実行")

	return sendCmd
}

func NewConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "設定管理",
		Long: `Webhook URLなどの設定を管理します。

利用可能な設定項目:
  webhook_url    Discord Webhook URL
                 形式: https://discord.com/api/webhooks/{id}/{token}
                 用途: メッセージ送信先のDiscord Webhook URL

設定ファイル:
  デフォルト: ~/.discord-webhook/config.json
  カスタム: --config フラグで指定可能

使用例:
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
		Short: "設定値を設定",
		Long:  "設定値を設定します。現在はwebhook_urlのみサポートしています。",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSet(configPath, args[0], args[1])
		},
	}

	setCmd.Flags().StringVarP(&configPath, "config", "c", "", "設定ファイルのパス")

	return setCmd
}

func NewConfigGetCommand() *cobra.Command {
	var configPath string

	getCmd := &cobra.Command{
		Use:   "get [key]",
		Short: "設定値を取得",
		Long:  "設定値を表示します。keyを指定しない場合はすべての設定を表示します。",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := ""
			if len(args) > 0 {
				key = args[0]
			}
			return runConfigGet(configPath, key)
		},
	}

	getCmd.Flags().StringVarP(&configPath, "config", "c", "", "設定ファイルのパス")

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