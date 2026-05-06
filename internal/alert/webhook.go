package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to a webhook endpoint.
type WebhookPayload struct {
	Alerts    []Alert   `json:"alerts"`
	SentAt    time.Time `json:"sent_at"`
	TotalCrit int       `json:"total_crit"`
	TotalWarn int       `json:"total_warn"`
}

// WebhookSender posts alerts to an HTTP endpoint.
type WebhookSender struct {
	URL    string
	client *http.Client
}

// NewWebhookSender creates a WebhookSender targeting the given URL.
func NewWebhookSender(url string, timeout time.Duration) *WebhookSender {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &WebhookSender{
		URL:    url,
		client: &http.Client{Timeout: timeout},
	}
}

// Send serialises alerts and POSTs them to the configured URL.
func (w *WebhookSender) Send(alerts []Alert) error {
	if len(alerts) == 0 {
		return nil
	}
	payload := WebhookPayload{
		Alerts: alerts,
		SentAt: time.Now().UTC(),
	}
	for _, a := range alerts {
		switch a.Level {
		case LevelCrit:
			payload.TotalCrit++
		case LevelWarn:
			payload.TotalWarn++
		}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("alert: marshal payload: %w", err)
	}
	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("alert: webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("alert: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
