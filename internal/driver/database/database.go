package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aidapedia/gdk/log"
	"github.com/aidapedia/jabberwock/pkg/config"
	"go.uber.org/zap"
)

var DatabaseDriver *sql.DB

func NewDatabase(ctx context.Context) *sql.DB {
	cfg := config.GetConfig(ctx)
	if cfg == nil {
		log.FatalCtx(ctx, "failed to connect database: %v", zap.Error(errors.New("config is nil")))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Secret.Database.Host, cfg.Secret.Database.Port, cfg.Secret.Database.Username, cfg.Secret.Database.Password, cfg.Secret.Database.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.FatalCtx(ctx, "failed to connect database: %v", zap.Error(err))
	}

	err = db.Ping()
	if err != nil {
		log.FatalCtx(ctx, "failed to connect database: %v", zap.Error(err))
	}

	DatabaseDriver = db
	return db
}

func BeginTx(ctx context.Context) (*sql.Tx, error) {
	return DatabaseDriver.BeginTx(ctx, nil)
}
