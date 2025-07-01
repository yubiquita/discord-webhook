# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

このファイルは、このリポジトリでコードを扱う際のClaude Code (claude.ai/code) 向けガイダンスを提供します。

## プロジェクト概要

これはGoで構築されたDiscord Webhook CLIツールで、webhook URLを介してDiscordチャンネルにメッセージを送信できます。このツールはwebhook URLの設定管理機能を提供し、一回限りのURL指定、永続的な設定保存、標準入力からのパイプ入力の全てをサポートしています。

**重要なモジュール情報**: `github.com/yubiquita/discord-webhook` - go installでの公開リポジトリ（v0.1.0でリリース済み）

## アーキテクチャ

コードベースは明確な関心の分離を持つクリーンアーキテクチャパターンに従っています：

- **`cmd/discord-webhook/main.go`**: CLIを初期化するアプリケーションエントリーポイント
- **`internal/cli/`**: Cobraフレームワークを使用したCLIコマンド定義とルーティング
  - `root.go`: 全てのCLIコマンドとそのフラグを定義
  - `commands.go`: 各コマンドの実際のビジネスロジックを含む（RunSend, RunConfigSet, RunConfigGet）
- **`internal/webhook/`**: Discord webhook API通信用HTTPクライアント（DiscordMessage構造体、Client構造体）
- **`internal/config/`**: JSON永続化による設定ファイル管理（`~/.discord-webhook/config.json`）

CLI層は`commands.go`のビジネスロジック関数に委譲し、configとwebhookパッケージ間を調整します。Cobraコマンドの各層は薄いラッパーとして機能し、テスト可能なビジネスロジック関数を呼び出します。

**重要な設計パターン**:
- `RunSend`関数は`cobra.Command`を受け取り、`cmd.InOrStdin()`と`cmd.OutOrStdout()`を使用してテスト可能な入出力を実現
- 標準入力からのメッセージ読み取りは`--message`フラグが未指定の場合のみ実行され、フラグが優先される
- webhook.Client構造体はHTTPクライアントをラップし、JSON形式でDiscord APIと通信

## 開発コマンド

### ビルド・インストール
```bash
# ローカルビルド
go build -o discord-webhook ./cmd/discord-webhook/

# グローバルインストール（開発用）
go install ./cmd/discord-webhook/

# クロスプラットフォームビルド（リリース用）
GOOS=linux GOARCH=amd64 go build -o discord-webhook-linux ./cmd/discord-webhook/
GOOS=windows GOARCH=amd64 go build -o discord-webhook.exe ./cmd/discord-webhook/
GOOS=darwin GOARCH=amd64 go build -o discord-webhook-darwin ./cmd/discord-webhook/
```

### テスト
```bash
# 全てのテストを実行
go test ./...

# 特定のパッケージのテストを実行
go test ./internal/webhook/
go test ./internal/config/
go test ./internal/cli/

# 詳細出力でテストを実行
go test -v ./...

# 特定のテストケースのみ実行（日本語テスト名対応）
go test ./internal/cli/ -run "TestSendCommand_標準入力からメッセージを読み取る"

# カバレッジ付きテスト
go test -cover ./...
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

### 依存関係
```bash
# 新しい依存関係を追加
go get <package>

# 依存関係を整理・検証
go mod tidy
go mod verify

# 脆弱性チェック
go list -json -m all | nancy sleuth
```

### リリース・タグ管理
```bash
# 新しいバージョンタグを作成（semver形式）
git tag v0.2.0
git push origin v0.2.0

# タグ一覧を確認
git tag -l

# 公開されたバージョンからインストール
go install github.com/yubiquita/discord-webhook/cmd/discord-webhook@latest
go install github.com/yubiquita/discord-webhook/cmd/discord-webhook@v0.1.0
```

## テストアーキテクチャ

このプロジェクトは包括的なテストカバレッジを持つTDD原則に従っています：

- **Unit tests**: 各パッケージには対応する`*_test.go`ファイルがある
- **Mock servers**: HTTPテストは`httptest.NewServer`を使用してDiscord APIレスポンスをモック
- **Japanese test names**: テスト関数は`TestFunction_Scenario_ExpectedBehavior`パターンに従った説明的な日本語名を使用
- **Table-driven tests**: 現在は使用されていないが、バリデーションテストを拡張する際に適切

## CLI コマンド構造

Cobraフレームワークを使用して以下の階層でコマンドが構築されています：
- `discord-webhook` (root)
  - `send` - Discord webhookにメッセージを送信
  - `config` - 設定管理
    - `set <key> <value>` - 設定値を設定
    - `get [key]` - 設定値を取得

`send`コマンドは複数の入力方法をサポートしています：
- インラインURL指定（`-u/--url`）
- 永続的な設定ファイル使用
- `--message`フラグによる直接指定
- 標準入力からのパイプ入力（`--message`未指定時）

`--dry-run`フラグは実際のメッセージ送信なしでテストを可能にします。

## Claude Code連携

このツールはClaude Codeのnotification hookとして使用可能です。`~/.claude/settings.json`での設定例：

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

## 使用例

```bash
# フラグでメッセージを指定
discord-webhook send --message "Hello" --url https://discord.com/api/webhooks/...

# 標準入力からメッセージを受け取り
echo "Hello Discord!" | discord-webhook send --url https://discord.com/api/webhooks/...

# ファイルからメッセージを送信
cat message.txt | discord-webhook send

# 設定ファイルを使用してドライラン
echo "Test message" | discord-webhook send --dry-run

# JSON入力をパイプで処理（Claude Code連携）
echo '{"message": "Build completed"}' | jq -r '.message' | discord-webhook send
```