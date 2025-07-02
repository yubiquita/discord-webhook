package webhook

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendMessage_ReturnsSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method: POST, actual: %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json, actual: %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient()
	err := client.SendMessage(server.URL, "test message")

	if err != nil {
		t.Errorf("error occurred: %v", err)
	}
}

func TestClient_SendMessage_ReturnsErrorForInvalidURL(t *testing.T) {
	client := NewClient()
	err := client.SendMessage("invalid-url", "test message")

	if err == nil {
		t.Error("no error returned for invalid URL")
	}
}

func TestClient_SendMessage_ReturnsErrorForEmptyMessage(t *testing.T) {
	client := NewClient()
	err := client.SendMessage("https://discord.com/api/webhooks/test", "")

	if err == nil {
		t.Error("no error returned for empty message")
	}
}

func TestClient_SendMessage_HandlesServerErrorProperly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient()
	err := client.SendMessage(server.URL, "test message")

	if err == nil {
		t.Error("no error returned for server error")
	}
}