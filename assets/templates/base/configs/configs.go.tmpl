package configs

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Provide a default config
var DefaultConfig = []byte(`
environment: D
grpc_address: 10443
http_address: 10080
max_concurrent_streams: 1000
`)

// Config holds all settings
type Config struct {
	HTTPAddress          int    `yaml:"http_address" mapstructure:"http_address"`
	GRPCAddress          int    `yaml:"grpc_address" mapstructure:"grpc_address"`
	Environment          string `yaml:"environment" mapstructure:"environment"`
	MaxConcurrentStreams uint32 `yaml:"max_concurrent_streams" mapstructure:"max_concurrent_streams"`
}

func Load(cfg *Config, defaultConfig []byte) error {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		return fmt.Errorf("failed to read viper configs %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal configs %w", err)
	}

	return nil
}
