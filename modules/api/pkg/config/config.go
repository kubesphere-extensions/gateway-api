package config

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/spf13/viper"

	gatewayapi "github.com/kubesphere-extensions/gateway-api/pkg/simple/gatewayapi/options"
)

// Package config saves configuration for running KubeSphere components
//
// Config can be configured from command line flags and configuration file.
// Command line flags hold higher priority than configuration file. But if
// component Endpoint/Host/APIServer was left empty, all of that component
// command line flags will be ignored, use configuration file instead.
// For example, we have configuration file
//
// gatewayapi:
//   namespace: kubesphere-controls-system
//   repository: kubesphere/nginx-ingress-controller
//   tag: v1.3.1
//   watchesPath: /var/helm-charts/watches.yaml
//
// At the same time, have command line flags like following:
// --gatewayapi-namespace kubesphere-system --gatewayapi-watchesPath ./watches.yaml
//
//  Command line has higher priority. But if command line flags like following:
// --gatewayapi-namespace ns --gatewayapi-watchesPath /watches.yaml

var (
	// singleton instance of config package
	_config = defaultConfig()
)

const (
	// DefaultConfigurationName is the default name of configuration
	defaultConfigurationName = "config"

	// DefaultConfigurationPath the default location of the configuration file
	defaultConfigurationPath = "/etc/gatewayapi"
)

type config struct {
	cfg         *Config
	cfgChangeCh chan Config
	loadOnce    sync.Once
}

func (c *config) loadFromDisk() (*Config, error) {
	var err error
	c.loadOnce.Do(func() {
		if err = viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				err = fmt.Errorf("error parsing configuration file %s", err)
			}
		}
		err = viper.Unmarshal(c.cfg)
	})
	return c.cfg, err
}

func defaultConfig() *config {
	viper.SetConfigName(defaultConfigurationName)
	viper.AddConfigPath(defaultConfigurationPath)

	// Load from current working directory, only used for debugging
	viper.AddConfigPath(".")

	// Load from Environment variables
	viper.SetEnvPrefix("gatewayapi")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &config{
		cfg:         New(),
		cfgChangeCh: make(chan Config),
		loadOnce:    sync.Once{},
	}
}

// Config defines everything needed for apiserver to deal with external services
type Config struct {
	GatewayOptions *gatewayapi.Options `json:"gatewayapi,omitempty" yaml:"gatewayapi,omitempty" mapstructure:"gatewayapi"`
}

// newConfig creates a default non-empty Config
func New() *Config {
	return &Config{
		GatewayOptions: gatewayapi.NewGatewayApiOptions(),
	}
}

// TryLoadFromDisk loads configuration from default location after server startup
// return nil error if configuration file not exists
func TryLoadFromDisk() (*Config, error) {
	return _config.loadFromDisk()
}

// convertToMap simply converts config to map[string]bool
// to hide sensitive information
func (conf *Config) ToMap() map[string]bool {

	result := make(map[string]bool, 0)

	if conf == nil {
		return result
	}

	c := reflect.Indirect(reflect.ValueOf(conf))

	for i := 0; i < c.NumField(); i++ {
		name := strings.Split(c.Type().Field(i).Tag.Get("json"), ",")[0]
		if strings.HasPrefix(name, "-") {
			continue
		}

		if c.Field(i).IsNil() {
			result[name] = false
		} else {
			result[name] = true
		}
	}

	return result
}
