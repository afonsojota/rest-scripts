package config

import (
	"github.com/afonsojota/go-afonsojota-toolkit/goutils/logger"
	"github.com/spf13/viper"
	"os"
	"path"
)

const (
	localScope          = "local"
	cleancoderDevScope        = "dev"
	cleancoderBetaScope       = "beta"
	cleancoderProductionScope = "prod"
)

var scope string

func GetScope() string {
	if scope != "" {
		return scope
	} else {
		envScope := os.Getenv("SCOPE")
		if envScope == "" || (envScope != cleancoderDevScope && envScope != cleancoderBetaScope && envScope != cleancoderProductionScope) {
			envScope = localScope
		}

		scope = envScope
		return scope
	}
}

func InLocal() bool {
	return localScope == GetScope()
}

func InDev() bool {
	return cleancoderDevScope == GetScope()
}

func InBeta() bool {
	return cleancoderBetaScope == GetScope()
}

func InProd() bool {
	return cleancoderProductionScope == GetScope()
}

func getConfigFile() string {
	if InProd() {
		return "config-prod.toml"
	} else if InDev() {
		return "config-dev.toml"
	} else if InBeta() {
		return "config-beta.toml"
	}
	return "config-local.toml"
}

func Init() {
	file := getConfigFile()
	path := getcleancoderPath(file)
	viper.SetConfigFile(path)
	logger.Debugf("Viper loaded configuration file: %s", path)
	logger.Debugf("Scope: %s", scope)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func getcleancoderPath(name string) string {
	confDir := os.Getenv("CONF_DIR")

	if confDir == "" {
		logger.Info("CONF_DIR env variable was not set")
		base := os.Getenv("GOPATH")
		appPath := os.Getenv("APPLICATION_REPO")
		if len(appPath) == 0 {
			appPath = "github.com/afonsojota/cleancoder_rest-scripts"
		}
		confDir = path.Join(base, "src", appPath, "conf")
		logger.Infof("Using %s as CONF_DIR", confDir)
	}
	logger.Infof("reading conf from: %s", confDir)

	confDir = confDir + "/" + name
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		panic("Not able to find configuration dir. Please check your CONF_DIR environment variable")
	}
	return confDir
}
