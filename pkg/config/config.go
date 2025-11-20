package config

import (
	"context"

	"os"

	gconfig "github.com/aidapedia/gdk/config"
	genv "github.com/aidapedia/gdk/environment"
	gfeatureflag "github.com/aidapedia/gdk/featureflag"
	gfeatureflagmodule "github.com/aidapedia/gdk/featureflag/module"
	"github.com/aidapedia/gdk/log"
	"github.com/aidapedia/gdk/util"
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
	var (
		value interface{}
		err   error
	)
	// Update Feature Flag
	// Add new feature flag here
	value, err = ff.GetValue(ctx, "disable_tracer")
	valBool := util.ToBool(value)
	if err != nil {
		log.WarnCtx(ctx, "Error updating feature flag", zap.String("key", "disable_tracer"), zap.Error(err))
	} else if valBool != globalConfig.FeatureFlags.DisableTracer {
		globalConfig.FeatureFlags.DisableTracer = valBool
	}
}
