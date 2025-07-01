package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
}

type DiscordMessage struct {
	Content string `json:"content"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SendMessage(webhookURL, content string) error {
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("メッセージが空です")
	}

	message := DiscordMessage{
		Content: content,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSONエンコードに失敗しました: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTPリクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Discord API エラー: ステータスコード %d", resp.StatusCode)
	}

	return nil
}