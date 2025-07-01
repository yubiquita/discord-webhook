# Discord Webhook CLI

Discord Webhook URLを使用してDiscordチャンネルにメッセージを送信するためのシンプルで強力なコマンドラインツールです。

## 特徴

- **複数の入力方法**: フラグ、標準入力、またはパイプからメッセージを受け取り
- **設定管理**: Webhook URLを永続的に保存して再利用
- **ドライラン**: 実際にメッセージを送信せずにテスト
- **クリーンな出力**: 人間が読みやすい形式でレスポンスを表示
- **Claude Code連携**: notification hookとして使用可能

## インストール

### go installでインストール（推奨）

```bash
go install github.com/yubiquita/discord-webhook/cmd/discord-webhook@latest
```

### ソースからビルド

```bash
git clone https://github.com/yubiquita/discord-webhook.git
cd discord-webhook
go build -o discord-webhook ./cmd/discord-webhook/
```

ビルドされたバイナリを任意の場所に配置してください（例：`/usr/local/bin/`）。

## 必要要件

- Go 1.24.4以上

## 使用方法

### 基本的な使用例

#### 1. フラグでメッセージを指定

```bash
discord-webhook send --message "Hello Discord!" --url https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

#### 2. 標準入力からメッセージを読み取り

```bash
echo "Hello from command line!" | discord-webhook send --url https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

#### 3. ファイルからメッセージを送信

```bash
cat message.txt | discord-webhook send --url https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

### 設定管理

Webhook URLを永続的に保存して、毎回指定する必要をなくすことができます。

#### Webhook URLを設定

```bash
discord-webhook config set webhook_url https://discord.com/api/webhooks/YOUR_WEBHOOK_URL
```

#### 設定されたURLを使用してメッセージを送信

```bash
discord-webhook send --message "設定されたURLを使用"
```

または

```bash
echo "パイプからのメッセージ" | discord-webhook send
```

#### 設定値を確認

```bash
# 全ての設定を表示
discord-webhook config get

# 特定の設定を表示
discord-webhook config get webhook_url
```

### ドライラン

実際にメッセージを送信せずにテストできます：

```bash
echo "テストメッセージ" | discord-webhook send --dry-run
```

### Claude Code notification hookとの連携

Claude Codeのnotification hook機能と組み合わせることで、様々なイベント（ファイル変更、コマンド実行など）をDiscordに通知できます。

#### 設定例

Claude Codeの設定ファイル（`~/.claude/settings.json`）でnotification hookを設定：

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

この設定により、Claude CodeからのJSON形式の通知メッセージを自動的にDiscordに転送できます。

## コマンドリファレンス

### `discord-webhook send`

Discordチャンネルにメッセージを送信します。

**フラグ:**
- `-m, --message string`: 送信するメッセージ（指定されない場合は標準入力から読み取り）
- `-u, --url string`: Discord Webhook URL（設定ファイルよりも優先）
- `--dry-run`: 実際のメッセージ送信なしでテスト

### `discord-webhook config`

設定を管理します。

**サブコマンド:**
- `set <key> <value>`: 設定値を設定
- `get [key]`: 設定値を取得（keyが未指定の場合は全設定を表示）

**設定可能なキー:**
- `webhook_url`: Discord Webhook URL

## 設定ファイル

設定は以下の場所に保存されます：
- `~/.discord-webhook/config.json`

ファイルはJSON形式で、以下のような構造です：

```json
{
  "webhook_url": "https://discord.com/api/webhooks/YOUR_WEBHOOK_URL"
}
```

## Discord Webhook URLの取得方法

1. Discordでメッセージを送信したいチャンネルを選択
2. チャンネル設定（歯車アイコン）を開く
3. 「連携サービス」→「ウェブフック」を選択
4. 「新しいウェブフック」を作成
5. 生成されたWebhook URLをコピー

## トラブルシューティング

### エラー: "No webhook URL provided"

- `--url`フラグでURLを指定するか、`config set webhook_url`で設定してください

### エラー: "Failed to send message"

- Webhook URLが正しいか確認してください
- インターネット接続を確認してください
- Discord側でWebhookが有効か確認してください

## 開発者向け情報

このプロジェクトはGoで開発されており、以下の主要な依存関係を使用しています：

- [Cobra](https://github.com/spf13/cobra): CLIフレームワーク
- 標準ライブラリのみでHTTPクライアントを実装

### テスト実行

```bash
go test ./...
```

### プロジェクト構造

```
.
├── cmd/discord-webhook/    # アプリケーションエントリーポイント
├── internal/cli/          # CLIコマンド定義
├── internal/webhook/      # Discord webhook API通信
└── internal/config/       # 設定ファイル管理
```

## ライセンス

このプロジェクトはMIT Licenseの下で公開されています。詳細は[LICENSE](LICENSE)ファイルを参照してください。