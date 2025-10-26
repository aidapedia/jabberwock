package main

import (
	"context"
	"os"
	"time"

	"github.com/aidapedia/gdk/log"
	"github.com/aidapedia/gdk/telemetry/tracer"
	"github.com/aidapedia/jabberwock/internal/app"
	"github.com/aidapedia/jabberwock/pkg/config"
	"github.com/aidapedia/jabberwock/pkg/jwt"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Initialize config
	cfg := config.GetConfig(ctx)

	// Initialize logger
	log.New(&log.Config{
		Level: log.LoggerLevel(cfg.App.Log.Level),
	})
	defer log.Sync()

	// Initialize JWT Token
	jwt.Init(cfg.Secret.Auth.PrivateKey)

	// Initialize tracer
	tr, err := tracer.InitTracer(cfg.App.Name, cfg.Secret.Tracer.AddressURL, false)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to initialize tracer", zap.Error(err))
		os.Exit(0)
	}
	defer tr.Shutdown(ctx)

	// Set timezone
	loc, err := time.LoadLocation(cfg.App.LocalTime)
	if err != nil {
		loc = time.Local
	}
	time.Local = loc

	// Initialize HTTP server
	err = app.InitHTTPServer().ListenAndServe()
	if err != nil {
		log.ErrorCtx(ctx, "Failed to initialize tracer", zap.Error(err))
		os.Exit(0)
	}
}
