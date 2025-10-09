package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/aidapedia/gdk/config"
	"github.com/aidapedia/gdk/environment"
	"github.com/aidapedia/jabberwock/pkg/config/model"
	"github.com/aidapedia/jabberwock/pkg/config/secret"
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
		"domea.main",
		"domea.storage",
		"domea.vendor",
	}
	configPaths = []string{
		"../../files/etc/domea", // for service configuration
		"/etc/domea",            // for service configuration in staging/production
		"files/etc/domea",       // for database migration
	}
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
	for _, conf := range configInstances {
		err := config.NewConfig(conf.files, conf.key, conf.config).SetConfig()
		if err != nil {
			return err
		}
	}
	// set secret config
	secretConfig, err := secret.GetSecret()
	if err != nil {
		return fmt.Errorf("failed to get secret: %s", err.Error())
	}

	globalConfig.Storage.PostgreSQL.Address = secretConfig.Database.Address
	globalConfig.Storage.PostgreSQL.Port = secretConfig.Database.Port
	globalConfig.Storage.PostgreSQL.Username = secretConfig.Database.Username
	globalConfig.Storage.PostgreSQL.Password = secretConfig.Database.Password
	globalConfig.Storage.PostgreSQL.Database = secretConfig.Database.Name
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
		Files:    getConfigFiles(files),
	}
}

func getConfigFiles(files []string) []string {
	result := make([]string, 0, len(files))
	for i := range files {
		result = append(result, fmt.Sprintf("%s.%s", files[i], environment.GetAppEnvironment()))
	}
	return result
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
