package config

import "github.com/spf13/viper"

// Config represents the structure of the configuration file
// Config represents the structure of the configuration file
type Config struct {
	BatchSize               int      `mapstructure:"batchSize"`
	NumGoroutinesMultiplier int      `mapstructure:"numGoroutinesMultiplier"`
	URLs                    []string `mapstructure:"urls"`
	FileURL                 string   `mapstructure:"fileURL"`
	OutputPath              string   `mapstructure:"outputPath"`
}

func LoadConfig(filePath string) (Config, error) {
	var config Config

	// Set the configuration file path
	viper.SetConfigFile(filePath)

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
