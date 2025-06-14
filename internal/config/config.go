package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Blockchain BlockchainConfig `mapstructure:"blockchain"`
	Network    NetworkConfig    `mapstructure:"network"`
	Miner      MinerConfig      `mapstructure:"miner"`
}

// BlockchainConfig holds blockchain-specific configuration
type BlockchainConfig struct {
	DifficultyCalculationBlocks int `mapstructure:"difficulty_calculation_blocks"`
	TargetBlockTime             int `mapstructure:"target_block_time"`
}

// NetworkConfig holds network-specific configuration
type NetworkConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// MinerConfig holds miner-specific configuration
type MinerConfig struct {
	NetworkSyncInterval int `mapstructure:"network_sync_interval"`
	MaxNonce            int `mapstructure:"max_nonce"`
}

// Load reads configuration from file
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// Default returns default configuration
func Default() *Config {
	return &Config{
		Blockchain: BlockchainConfig{
			DifficultyCalculationBlocks: 50,
			TargetBlockTime:             20,
		},
		Network: NetworkConfig{
			Host: "127.0.0.1",
			Port: 8080,
		},
		Miner: MinerConfig{
			NetworkSyncInterval: 1,
			MaxNonce:            4294967296,
		},
	}
}
