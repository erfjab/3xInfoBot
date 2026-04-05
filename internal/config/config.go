package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`
}

var Cfg *Config

func LoadConfig() (*Config, error) {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		var configNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configNotFound) {
			return &config, err
		}
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return &config, err
	}

	if err = config.validate(); err != nil {
		return &config, err
	}

	Cfg = &config
	return &config, nil
}

func (c *Config) validate() error {
	if c.TelegramBotToken == "" {
		return errors.New("TELEGRAM_BOT_TOKEN cannot be empty")
	}

	return nil
}