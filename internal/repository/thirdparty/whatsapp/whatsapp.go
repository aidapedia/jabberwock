package whatsapp

import (
	"context"
	"fmt"

	gcallwrapper "github.com/aidapedia/gdk/callwrapper"
	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/bytedance/sonic"
	"github.com/kurniajigunawan/homestay/internal/common/callwrapper"
)

// SendMessageText sends a text message to the specified phone number.
func (w *Whatsapp) SendMessageText(ctx context.Context, req SendTextRequest) (resp SendTextResponse, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "WhatsappThirdParty/SendMessageText")
	defer span.Finish(err)

	_, err = gcallwrapper.GetCallWrapper(callwrapper.WhatsappSendMessageText).Call(ctx, map[string]interface{}{
		"chatID": req.ChatID,
	}, func() (interface{}, error) {
		httpRes, err := w.httpClient.Send(ctx, req.ToHTTPRequest(ctx))
		if err != nil {
			return resp, err
		}
		defer httpRes.Close()

		if httpRes.StatusCode() >= 300 {
			return resp, fmt.Errorf("unexpected status code: %d", httpRes.StatusCode())
		}

		err = sonic.Unmarshal(httpRes.Body(), &resp)
		if err != nil {
			return resp, err
		}

		return resp, nil
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
