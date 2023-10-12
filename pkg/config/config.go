package config

import "github.com/spf13/viper"

type Config struct {
	BackendPort  string `mapstructure:"backend-port"`
	FrontendPort string `mapstructure:"frontend-port"`
}

var config *Config

func Init() (*Config, error) {

	viper.AddConfigPath("configs/")
	viper.SetConfigName("local")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	// If struct Config contains embedded structs, you have to use UnmarshalKey
	// and evidently write this keys and embedded structs
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	config = &cfg
	return &cfg, nil
}

func GetConfig() *Config {
	return config
}
