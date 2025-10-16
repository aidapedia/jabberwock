package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/aidapedia/gdk/config"
	gsecret "github.com/aidapedia/gdk/config/secret"
	"github.com/aidapedia/gdk/environment"
	"github.com/aidapedia/jabberwock/pkg/config/model"
)

// serviceConfig name of service config
const serviceConfig = "ServiceConfig"

// Configuration Instances
type configurationInstance struct {
	files  []config.FileConfig
	key    string
	config interface{}
}

var (
	globalConfig = &model.ServiceConfig{}
	doOnce       = sync.Once{}
)

var (
	configFiles = []string{
		"main",
		"storage",
		// "vendor",
	}
	configPaths = []string{
		"../../files/config", // for service configuration
		"/config",            // for service configuration in staging/production
		"files/config",       // for database migration
	}
	secretPath = "../../files/config/secret.json"
)

func GetConfig() *model.ServiceConfig {
	doOnce.Do(func() {
		if err := setInstances(); err != nil {
			log.Fatalf("Error setting config: %v", err)
		}
	})
	return getInstance()
}

func setConfig(configInstances []configurationInstance) error {
	var err error
	// cold config
	for _, conf := range configInstances {
		err = config.NewConfig(conf.files, conf.key, conf.config).SetConfig()
		if err != nil {
			return err
		}
	}
	// set secret config
	var secret gsecret.Vault
	if environment.GetAppEnvironment() == environment.Development {
		secret = gsecret.NewSecretFile(secretPath)
	}
	err = secret.GetSecret(context.Background(), &globalConfig.Secret)
	if err != nil {
		return fmt.Errorf("failed to get secret: %s", err.Error())
	}
	return nil
}

func setInstances() error {
	var configs []config.FileConfig
	for i := range configPaths {
		cfg := getFilesConfiguration(configPaths[i], configFiles)
		if _, err := os.Stat(cfg.FilePath); os.IsNotExist(err) {
			continue
		}
		configs = append([]config.FileConfig{}, cfg)
	}
	if len(configs) == 0 {
		return fmt.Errorf("no configuration file found")
	}

	configInstances := getConfigurationInstances(configs)
	err := setConfig(configInstances)
	if err != nil {
		return err
	}
	return err
}

func getFilesConfiguration(path string, files []string) config.FileConfig {
	return config.FileConfig{
		FilePath: filepath.Join(path, environment.GetAppEnvironment()),
		Files:    files,
	}
}

func getConfigurationInstances(cfg []config.FileConfig) []configurationInstance {
	return []configurationInstance{
		{
			files:  cfg,
			key:    serviceConfig,
			config: globalConfig,
		},
	}
}

func getInstance() *model.ServiceConfig {
	return globalConfig
}
