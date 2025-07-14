package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var ConfigInstance *Config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Security SecurityConfig
	Redis    RedisConfig
}

type RateLimitConfig struct {
	RequestPerSec int
	Burst         int
}

type ServerConfig struct {
	Host      string
	PortGrpc  string
	PortHttp  string
	ApiKey    string
	RateLimit RateLimitConfig
}

type SecurityConfig struct {
	SecretKey  string
	Expiration int64
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Đọc tệp cấu hình
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	// Ánh xạ dữ liệu YAML vào struct
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Lưu cấu hình vào biến toàn cục
	ConfigInstance = &config

	return config, nil
}

func GetConfig() *Config {
	if ConfigInstance == nil {
		log.Fatal("ConfigInstance is not initialized. Please call LoadConfig first.")
	}
	return ConfigInstance
}
