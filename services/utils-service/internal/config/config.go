package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host                  string        `mapstructure:"SERVER_HOST"`
	Port                  string        `mapstructure:"SERVER_PORT"`
	TokenSymmetricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	CloudinaryURL         string        `mapstructure:"CLOUDINARY_URL"`
	GoogleCredentialsPath string        `mapstructure:"GOOGLE_CREDENTIALS_PATH"`
	GoogleDriveFolderID   string        `mapstructure:"GOOGLE_DRIVE_FOLDER_ID"`
	StorageProvider       StorageProvider
}

func LoadConfig(storageProvider StorageProvider) (config Config, err error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config: %w", err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		return config, err
	}

	// assign from main
	config.StorageProvider = storageProvider

	// validate
	if !storageProvider.IsValid() {
		return config, fmt.Errorf("invalid storage provider: %s", storageProvider)
	}

	return config, nil
}
