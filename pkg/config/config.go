package config

import (
	"context"
	"log"
	"sync"

	gconfig "github.com/aidapedia/gdk/config"
	genv "github.com/aidapedia/gdk/environment"
	"github.com/aidapedia/jabberwock/pkg/config/model"
)

var (
	globalConfig = &model.ServiceConfig{}
	secretConfig = &model.SecretConfig{}
	doOnce       = sync.Once{}
)

var configManager *gconfig.Manager

func init() {
	secretType := gconfig.SecretTypeFile
	if genv.GetAppEnvironment() != genv.Development {
		secretType = gconfig.SecretTypeGSM
	}
	configManager = gconfig.New(gconfig.Option{
		TargetStore: globalConfig,
		ConfigKey:   "ServiceConfig",
		FileName: []string{
			"main",
			"storage",
		},
		WithSecret:   secretType,
		TargetSecret: secretConfig,
	})
}

func GetConfig(ctx context.Context) *model.ServiceConfig {
	doOnce.Do(func() {
		if err := configManager.SetConfig(ctx); err != nil {
			log.Fatalf("Error setting config: %v", err)
		}
	})
	res, _ := configManager.GetConfig(ctx)
	cfg := res.(*model.ServiceConfig)
	cfg.Secret = *getSecret(ctx)
	return cfg
}

func getSecret(ctx context.Context) *model.SecretConfig {
	doOnce.Do(func() {
		if err := configManager.SetSecretStore(ctx); err != nil {
			log.Fatalf("Error setting config: %v", err)
		}
	})
	res, _ := configManager.GetSecret(ctx)
	return res.(*model.SecretConfig)
}
