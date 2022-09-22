package config

import (
	"github.com/spf13/viper"
)

var (
	Config *ConfigStruct
)

type ConfigStruct struct {
	// RPC Config
	GRPC_PORT int `mapstructure:"RPC_PORT"`
	// End of RPC Config

	// Database Config
	DB_USER        string `mapstructure:"DB_USER"`
	DB_PASSWORD    string `mapstructure:"DB_PASSWORD"`
	DB_HOST        string `mapstructure:"DB_HOST"`
	DB_PORT        int    `mapstructure:"DB_PORT"`
	DB_NAME        string `mapstructure:"DB_NAME"`
	REDIS_HOST     string `mapstructure:"REDIS_HOST"`
	REDIS_PORT     int    `mapstructure:"REDIS_PORT"`
	REDIS_USER     string `mapstructure:"REDIS_USER"`
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	// End of Database Config

	// Auth Config
	ACCESS_TOKEN_EXPIRED  int64  `mapstructure:"ACCESS_TOKEN_EXPIRED"`
	REFRESH_TOKEN_EXPIRED int64  `mapstructure:"REFRESH_TOKEN_EXPIRED"`
	JWT_SECRET            string `mapstructure:"JWT_SECRET"`
	// End of Auth Config

	// REST API Config
	REST_API_PORT        int      `mapstructure:"REST_API_PORT"`
	REST_API_HOST        string   `mapstructure:"REST_API_HOST"`
	CORS_ALLOWED_ORIGINS []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	// End of REST API Config
}

func LoadConfig() (config *ConfigStruct, err error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
