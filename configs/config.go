package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Biến toàn cục để lưu cấu hình
var ConfigInstance *Config

// Config chứa các giá trị cấu hình
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Security SecurityConfig
}

// ServerConfig chứa cấu hình server
type ServerConfig struct {
    Host string
    Port string
    ApiKey string 
}

type SecurityConfig struct {
    SecretKey string
}

// DatabaseConfig chứa cấu hình database
type DatabaseConfig struct {
    User     string
    Password string
    Host     string
    Port     string
    DbName   string
}

// LoadConfig khởi tạo và đọc cấu hình từ Viper
func LoadConfig() (config Config, err error) {
    // Thiết lập tên tệp cấu hình và định dạng
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
