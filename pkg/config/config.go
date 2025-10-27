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
	globalConfig = model.ServiceConfig{}
	doOnce       = sync.Once{}
)

var configManager *gconfig.Manager

func init() {
	secretType := gconfig.SecretTypeFile
	if genv.GetAppEnvironment() != genv.Development {
		secretType = gconfig.SecretTypeGSM
	}
	configManager = gconfig.New(gconfig.Option{
		TargetStore: &globalConfig,
		ConfigKey:   "ServiceConfig",
		FileName: []string{
			"main",
			"storage",
		},
		WithSecret:   secretType,
		TargetSecret: &globalConfig.Secret,
	})
	doOnce.Do(func() {
		if err := configManager.SetConfig(context.Background()); err != nil {
			log.Fatalf("Error setting config: %v", err)
		}
		if err := configManager.SetSecretStore(context.Background()); err != nil {
			log.Fatalf("Error setting config: %v", err)
		}
	})
}

func GetConfig(ctx context.Context) *model.ServiceConfig {
	return &globalConfig
}
