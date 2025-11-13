package main

import (
	"context"
	"os"
	"time"

	gjwt "github.com/aidapedia/gdk/cryptography/jwt"
	"github.com/aidapedia/gdk/log"
	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/aidapedia/jabberwock/internal/app"
	"github.com/aidapedia/jabberwock/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Initialize config
	cfg := config.GetConfig(ctx)

	// Initialize JWT
	gjwt.New([]byte(cfg.Secret.Auth.PrivateKey), []byte(cfg.Secret.Auth.PublicKey), []jwt.ParserOption{
		jwt.WithExpirationRequired(),
	}...)

	// Initialize logger
	log.New(&log.Config{
		Level:  log.LoggerLevel(cfg.App.Log.Level),
		Caller: false,
		DefaultTags: map[string]interface{}{
			"app": cfg.App.Name,
		},
	})
	defer log.Sync()

	// Initialize tracer
	if cfg.App.FeatureFlags.UseTracer {
		tr, err := tracer.InitTracer(cfg.App.Name, cfg.Secret.Tracer.AddressURL, false)
		if err != nil {
			log.ErrorCtx(ctx, "Failed to initialize tracer", zap.Error(err))
			os.Exit(0)
		}
		defer tr.Shutdown(ctx)
	}

	// Set timezone
	// This catch case like if your server on Singapore Datacenter but you wanna local time to UTC you have to set manually.
	loc, err := time.LoadLocation(cfg.App.LocalTime)
	if err != nil {
		loc = time.Local
	}
	time.Local = loc

	// Initialize HTTP server
	err = app.InitHTTPServer(ctx).ListenAndServe()
	if err != nil {
		log.ErrorCtx(ctx, "Failed to initialize tracer", zap.Error(err))
		os.Exit(0)
	}
}
