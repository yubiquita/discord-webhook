# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際のClaude Code (claude.ai/code) 向けガイダンスを提供します。

## プロジェクト概要

これはGoで構築されたDiscord Webhook CLIツールで、webhook URLを介してDiscordチャンネルにメッセージを送信できます。このツールはwebhook URLの設定管理機能を提供し、一回限りのURL指定と永続的な設定保存の両方をサポートしています。

## アーキテクチャ

コードベースは明確な関心の分離を持つクリーンアーキテクチャパターンに従っています：

- **`cmd/discord-webhook/main.go`**: CLIを初期化するアプリケーションエントリーポイント
- **`internal/cli/`**: Cobraフレームワークを使用したCLIコマンド定義とルーティング
  - `root.go`: 全てのCLIコマンドとそのフラグを定義
  - `commands.go`: 各コマンドの実際のビジネスロジックを含む
- **`internal/webhook/`**: Discord webhook API通信用HTTPクライアント
- **`internal/config/`**: JSON永続化による設定ファイル管理

CLI層は`commands.go`のビジネスロジック関数に委譲し、configとwebhookパッケージ間を調整します。設定はデフォルトで`~/.discord-webhook/config.json`にJSONとして保存されます。

## 開発コマンド

### ビルド
```bash
go build -o discord-webhook ./cmd/discord-webhook/
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
```

### 依存関係
```bash
# 新しい依存関係を追加
go get <package>

# 依存関係を整理
go mod tidy
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

`send`コマンドはインラインURL指定（`-u/--url`）と永続的な設定ファイル使用の両方をサポートしています。`--dry-run`フラグは実際のメッセージ送信なしでテストを可能にします。