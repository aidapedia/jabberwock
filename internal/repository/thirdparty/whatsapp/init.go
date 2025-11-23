package whatsapp

import (
	"context"
	"log"

	gcallwrapper "github.com/aidapedia/gdk/callwrapper"
	ghttpc "github.com/aidapedia/gdk/http/client"
	"github.com/kurniajigunawan/homestay/internal/common/callwrapper"
)

type Interface interface {
	SendMessageText(ctx context.Context, req SendTextRequest) (resp SendTextResponse, err error)
}

type Whatsapp struct {
	httpClient *ghttpc.Client
}

func New() Interface {
	// Init Callwrapper
	err := gcallwrapper.New(callwrapper.WhatsappSendMessageText, gcallwrapper.Options{Singleflight: true})
	if err != nil {
		log.Fatal(err)
	}
	gcallwrapper.GetCallWrapper(callwrapper.WhatsappSendMessageText)
	return &Whatsapp{
		httpClient: ghttpc.New(),
	}
}
