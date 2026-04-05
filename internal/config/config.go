package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken    string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	PanelURL            string `mapstructure:"PANEL_URL"`
	PanelPath           string `mapstructure:"PANEL_PATH"`
	PanelUsername       string `mapstructure:"PANEL_USERNAME"`
	PanelPassword       string `mapstructure:"PANEL_PASSWORD"`
	PanelTwoFactorCode  string `mapstructure:"PANEL_TWO_FACTOR_CODE"`
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
	if c.PanelURL == "" {
		return errors.New("PANEL_URL cannot be empty")
	}
	if c.PanelPath == "" {
		return errors.New("PANEL_PATH cannot be empty")
	}
	if c.PanelUsername == "" {
		return errors.New("PANEL_USERNAME cannot be empty")
	}
	if c.PanelPassword == "" {
		return errors.New("PANEL_PASSWORD cannot be empty")
	}

	return nil
}