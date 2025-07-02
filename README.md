# Discord Webhook CLI

A simple and powerful command-line tool for sending messages to Discord channels using Discord Webhook URLs.

## Features

- **Multiple input methods**: Accept messages from flags, stdin, or pipes
- **Configuration management**: Persistently save and reuse webhook URLs
- **Dry run**: Test without actually sending messages
- **Clean output**: Display responses in human-readable format
- **Claude Code integration**: Can be used as a notification hook

## Installation

### Install via go install (recommended)

```bash
go install github.com/yubiquita/discord-webhook/cmd/discord-webhook@latest
```

### Build from source

```bash
git clone https://github.com/yubiquita/discord-webhook.git
cd discord-webhook
go build -o discord-webhook ./cmd/discord-webhook/
```

Place the built binary in your desired location (e.g., `/usr/local/bin/`).

## Requirements

- Go 1.24.4 or higher

## Usage

### Basic usage examples

#### 1. Specify message with flag

```bash
discord-webhook send --message "Hello Discord!" --url https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

#### 2. Read message from stdin

```bash
echo "Hello from command line!" | discord-webhook send --url https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

#### 3. Send message from file

```bash
cat message.txt | discord-webhook send --url https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

### Configuration management

You can persistently save webhook URLs to avoid specifying them every time.

#### Set webhook URL

```bash
discord-webhook config set webhook_url https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

#### Send message using configured URL

```bash
discord-webhook send --message "Using configured URL"
```

or

```bash
echo "Message from pipe" | discord-webhook send
```

#### Check configuration values

```bash
# Display all configurations
discord-webhook config get

# Display specific configuration
discord-webhook config get webhook_url
```

### Dry run

You can test without actually sending messages:

```bash
echo "Test message" | discord-webhook send --dry-run
```

### Integration with Claude Code notification hooks

You can combine this tool with Claude Code's notification hook feature to send various events (file changes, command executions, etc.) to Discord.

#### Configuration example

Set up notification hooks in Claude Code's configuration file (`~/.claude/settings.json`):

```json
{
  "hooks": {
    "Notification": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "jq -r '.message' | discord-webhook send"
          }
        ]
      }
    ]
  }
}
```

This configuration automatically forwards JSON-formatted notification messages from Claude Code to Discord.

## Command reference

### `discord-webhook send`

Send a message to a Discord channel.

**Flags:**
- `-m, --message string`: Message to send (reads from stdin if not specified)
- `-u, --url string`: Discord Webhook URL (takes priority over configuration file)
- `--dry-run`: Test without actually sending the message

### `discord-webhook config`

Manage configurations.

**Subcommands:**
- `set <key> <value>`: Set configuration value
- `get [key]`: Get configuration value (displays all configurations if key is not specified)

**Configurable keys:**
- `webhook_url`: Discord Webhook URL

## Configuration file

Configurations are saved at:
- `~/.discord-webhook/config.json`

The file is in JSON format with the following structure:

```json
{
  "webhook_url": "https://discord.com/api/webhooks/YOUR_WEBHOOK_URL"
}
```

## How to get Discord Webhook URL

1. Select the Discord channel where you want to send messages
2. Open channel settings (gear icon)
3. Go to "Integrations" → "Webhooks"
4. Create a "New Webhook"
5. Copy the generated Webhook URL

## Troubleshooting

### Error: "No webhook URL provided"

- Specify a URL with the `--url` flag or set it with `config set webhook_url`

### Error: "Failed to send message"

- Verify that the webhook URL is correct
- Check your internet connection
- Confirm that the webhook is enabled on Discord

## Developer information

This project is developed in Go and uses the following main dependencies:

- [Cobra](https://github.com/spf13/cobra): CLI framework
- HTTP client implemented with standard library only

### Running tests

```bash
go test ./...
```

### Project structure

```
.
├── cmd/discord-webhook/    # Application entry point
├── internal/cli/          # CLI command definitions
├── internal/webhook/      # Discord webhook API communication
└── internal/config/       # Configuration file management
```

## License

This project is released under the MIT License. See the [LICENSE](LICENSE) file for details.