package whatsapp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	ghttpc "github.com/aidapedia/gdk/http/client"
	"github.com/kurniajigunawan/homestay/pkg/config"
)

type message string

func NewMessage(template Template, params ...interface{}) message {
	return message(fmt.Sprintf(string(template), params...))
}

type chatID string

// Create a phone number with format "6285111111111@c.us"
// Will normalize phone number
// Will add @c.us if phone number is not in format "6285111111111@c.us"
// Will remove + if phone number is in format "+6285111111111"
// Will remove space if phone number is in format "62 851 1111 1111"
func NewChatID(id string) chatID {
	id = strings.ReplaceAll(id, " ", "")
	if !strings.Contains(id, "@c.us") {
		id += "@c.us"
	}
	id = strings.TrimPrefix(id, "+")
	return chatID(id)
}

type SendTextRequest struct {
	ChatID  chatID  `json:"chatId"`
	Text    message `json:"text"`
	Session string  `json:"session"`
}

func (e SendTextRequest) ToHTTPRequest(ctx context.Context) *ghttpc.Request {
	cfg := config.GetConfig(ctx)
	httpReq := ghttpc.NewRequest(ctx)
	httpReq.SetMethod(http.MethodPost)
	httpReq.SetURL(cfg.Secret.Whatsapp.Host + "/api/sendText")
	httpReq.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"X-Api-Key":    cfg.Secret.Whatsapp.APIKey,
	})
	httpReq.SetJSON(e)
	return httpReq
}

type SendTextResponse struct {
	ID struct {
		FromMe     bool   `json:"fromMe"`
		Remote     string `json:"remote"`
		ID         string `json:"id"`
		Serialized string `json:"_serialized"`
	}
}
