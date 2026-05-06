package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend_EmptyAlerts_DoesNotPost(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	ws := NewWebhookSender(srv.URL, time.Second)
	if err := ws.Send(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty alerts")
	}
}

func TestSend_PostsPayload(t *testing.T) {
	var received WebhookPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ws := NewWebhookSender(srv.URL, time.Second)
	alerts := []Alert{
		{Host: "db01", Path: "/etc/mysql/my.cnf", Level: LevelCrit, Message: "4 line(s) changed"},
		{Host: "web01", Path: "/etc/nginx/nginx.conf", Level: LevelWarn, Message: "2 line(s) changed"},
	}
	if err := ws.Send(alerts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Alerts) != 2 {
		t.Errorf("expected 2 alerts in payload, got %d", len(received.Alerts))
	}
	if received.TotalCrit != 1 {
		t.Errorf("expected TotalCrit=1, got %d", received.TotalCrit)
	}
	if received.TotalWarn != 1 {
		t.Errorf("expected TotalWarn=1, got %d", received.TotalWarn)
	}
}

func TestSend_ServerError_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ws := NewWebhookSender(srv.URL, time.Second)
	alerts := []Alert{{Host: "h", Path: "/p", Level: LevelInfo, Message: "1 line(s) changed"}}
	if err := ws.Send(alerts); err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestNewWebhookSender_DefaultTimeout(t *testing.T) {
	ws := NewWebhookSender("http://example.com", 0)
	if ws.client.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %v", ws.client.Timeout)
	}
}
