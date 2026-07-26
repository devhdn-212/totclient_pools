package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/devhdn-212/totclient_api/internal/config"
)

// SendTelegramAlert notifies the ops Telegram chat about a server-side error.
// Fire-and-forget: callers run this in a goroutine so a slow/unreachable
// Telegram API never blocks the request that triggered the log. Failures are
// written straight to stderr (not through the zap logger) to avoid feeding
// back into the same error hook that calls this function.
func SendTelegramAlert(cfg config.Telegram, level, location, message string) {
	if cfg.BotToken == "" || cfg.ChatID == "" {
		return
	}

	text := fmt.Sprintf(
		"🚨 *%s* di `totclient_api`\n*Lokasi:* `%s`\n*Pesan:* %s",
		level, location, message,
	)

	payload, err := json.Marshal(map[string]string{
		"chat_id":    cfg.ChatID,
		"text":       text,
		"parse_mode": "Markdown",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "telegram alert: marshal error:", err)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	url := "https://api.telegram.org/bot" + cfg.BotToken + "/sendMessage"
	resp, err := client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintln(os.Stderr, "telegram alert: send error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		fmt.Fprintln(os.Stderr, "telegram alert: non-2xx response:", resp.StatusCode)
	}
}