package viper

import (
	"strings"

	"github.com/definancialbr/golang-container-kit/pkg/configuration"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	DefaultEnvKeyReplacer = strings.NewReplacer("-", "_", ".", "_")
)

type ConfigurationServiceOption func(*ConfigurationService)

type ConfigurationService struct {
	viper                    *viper.Viper
	optionalConfiguratioFile bool
}

func WithOptionalConfigurationFile() ConfigurationServiceOption {
	return func(c *ConfigurationService) {
		c.optionalConfiguratioFile = true
	}
}

func WithConfiguration(key string, defaultValue interface{}) ConfigurationServiceOption {
	return func(c *ConfigurationService) {
		c.viper.SetDefault(key, defaultValue)
	}
}

func WithFileType(fileType string) ConfigurationServiceOption {
	return func(c *ConfigurationService) {
		c.viper.SetConfigType(fileType)
	}
}

func WithFileName(fileName string) ConfigurationServiceOption {
	return func(c *ConfigurationService) {
		c.viper.SetConfigName(fileName)
	}
}

func WithEnvPrefix(envPrefix string) ConfigurationServiceOption {
	return func(c *ConfigurationService) {
		c.viper.SetEnvPrefix(envPrefix)
	}
}

func WithSearchPaths(searchPaths ...string) ConfigurationServiceOption {
	return func(c *ConfigurationService) {
		for _, searchPath := range searchPaths {
			c.viper.AddConfigPath(searchPath)
		}
	}
}

func WithHomeSearchPath() ConfigurationServiceOption {
	return func(c *ConfigurationService) {

		searchPath, err := homedir.Dir()
		if err != nil {
			return
		}

		searchPath, err = homedir.Expand(searchPath)
		if err != nil {
			return
		}

		c.viper.AddConfigPath(searchPath)

	}
}

func NewConfigurationService(options ...ConfigurationServiceOption) *ConfigurationService {

	c := &ConfigurationService{
		viper:                    viper.New(),
		optionalConfiguratioFile: false,
	}

	c.viper.SetConfigType("env")
	c.viper.SetEnvKeyReplacer(DefaultEnvKeyReplacer)
	c.viper.AutomaticEnv()

	for _, option := range options {
		option(c)
	}

	return c

}

func (c *ConfigurationService) Read() error {
	err := c.viper.ReadInConfig()

	if err != nil {

		_, ok := err.(viper.ConfigFileNotFoundError)

		if ok && c.optionalConfiguratioFile {
			return nil
		}
		return err

	}

	return nil
}

func (c *ConfigurationService) Load() configuration.Loader {
	return c.viper
}
