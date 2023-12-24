package config

import "github.com/spf13/viper"

type Config struct {
	BatchSize               int      `mapstructure:"batchSize"`
	NumGoroutinesMultiplier int      `mapstructure:"numGoroutinesMultiplier"`
	URLs                    []string `mapstructure:"urls"`
	FileURL                 string   `mapstructure:"fileURL"`
	OutputPath              string   `mapstructure:"outputPath"`
	TopNums                 int      `mapstructure:"TopNums"`
}

func LoadConfig(filePath string) (Config, error) {
	var config Config

	viper.SetConfigFile(filePath)

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
