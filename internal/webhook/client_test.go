package webhook

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendMessage_正常なレスポンスを返す(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("期待されるメソッド: POST, 実際: %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("期待されるContent-Type: application/json, 実際: %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient()
	err := client.SendMessage(server.URL, "テストメッセージ")

	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
}

func TestClient_SendMessage_不正なURLでエラーを返す(t *testing.T) {
	client := NewClient()
	err := client.SendMessage("invalid-url", "テストメッセージ")

	if err == nil {
		t.Error("不正なURLに対してエラーが返されませんでした")
	}
}

func TestClient_SendMessage_空のメッセージでエラーを返す(t *testing.T) {
	client := NewClient()
	err := client.SendMessage("https://discord.com/api/webhooks/test", "")

	if err == nil {
		t.Error("空のメッセージに対してエラーが返されませんでした")
	}
}

func TestClient_SendMessage_サーバーエラーを適切に処理する(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient()
	err := client.SendMessage(server.URL, "テストメッセージ")

	if err == nil {
		t.Error("サーバーエラーに対してエラーが返されませんでした")
	}
}