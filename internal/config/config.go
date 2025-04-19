package config

import (
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Solana struct {
		Endpoint string `mapstructure:"endpoint"`
		Network  string `mapstructure:"network"` // mainnet, testnet, devnet
	} `mapstructure:"solana"`

	Wallet struct {
		PrivateKey string `mapstructure:"private_key"`
	} `mapstructure:"wallet"`

	Monitoring struct {
		Wallets      []string `mapstructure:"wallets"`
		PollInterval int      `mapstructure:"poll_interval"` // in seconds
		ConfirmLevel int      `mapstructure:"confirm_level"`
	} `mapstructure:"monitoring"`

	Strategy struct {
		Type           string   `mapstructure:"type"` // mirror, filter, etc.
		MaxSlippage    float64  `mapstructure:"max_slippage"`
		MinTradeSize   float64  `mapstructure:"min_trade_size"`
		MaxTradeSize   float64  `mapstructure:"max_trade_size"`
		TradeDelay     int      `mapstructure:"trade_delay"` // in milliseconds
		TokenWhitelist []string `mapstructure:"token_whitelist"`
		TokenBlacklist []string `mapstructure:"token_blacklist"`
	} `mapstructure:"strategy"`
}

// Load loads the configuration from a file
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
