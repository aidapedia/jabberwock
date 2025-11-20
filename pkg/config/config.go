package config

import (
	"context"

	"os"

	gconfig "github.com/aidapedia/gdk/config"
	genv "github.com/aidapedia/gdk/environment"
	gfeatureflag "github.com/aidapedia/gdk/featureflag"
	gfeatureflagmodule "github.com/aidapedia/gdk/featureflag/module"
	"github.com/aidapedia/gdk/log"
	"github.com/aidapedia/jabberwock/pkg/config/model"
	"go.uber.org/zap"
)

var (
	globalConfig = model.ServiceConfig{}
)

var configManager *gconfig.Manager

func init() {
	ctx := context.Background()
	secretType := gconfig.SecretTypeFile
	if genv.GetAppEnvironment() != genv.Development {
		secretType = gconfig.SecretTypeGSM
	}
	configManager = gconfig.New(gconfig.Option{
		TargetStore: &globalConfig,
		ConfigKey:   "ServiceConfig",
		FileName: []string{
			"main",
		},
		WithSecret:   secretType,
		TargetSecret: &globalConfig.Secret,
	})
	if err := configManager.SetConfig(ctx); err != nil {
		log.FatalCtx(ctx, "Error setting config", zap.Error(err))
	}
	if err := configManager.SetSecretStore(ctx); err != nil {
		log.FatalCtx(ctx, "Error setting secret store", zap.Error(err))
	}
	// Initialize Feature Flag
	ffm := gfeatureflag.New(gfeatureflag.Option{
		Address: os.Getenv("FEATURE_FLAG_FILE_PATH"),
		Module:  gfeatureflagmodule.FileModule,
	})
	ffsignal, err := ffm.Watch(ctx)
	if err != nil {
		log.FatalCtx(ctx, "Error watching feature flag", zap.Error(err))
	}
	updateFeatureFlag(ctx, ffm)
	go func() {
		for range ffsignal {
			updateFeatureFlag(ctx, ffm)
		}
	}()
}

func GetConfig(ctx context.Context) *model.ServiceConfig {
	return &globalConfig
}

func updateFeatureFlag(ctx context.Context, ff gfeatureflagmodule.Interface) {
	// Update Feature Flag
	// Add new feature flag here
	key, err := ff.GetBool(ctx, "enable_tracer")
	if err != nil {
		log.WarnCtx(ctx, "Error updating feature flag", zap.String("key", "enable_tracer"), zap.Error(err))
	} else if key != globalConfig.FeatureFlags.UseTracer {
		globalConfig.FeatureFlags.UseTracer = key
	}
}
