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

    // Thiết lập giá trị mặc định
    viper.SetDefault("server.host", "127.0.0.1")
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("database.user", "12345678")
    viper.SetDefault("database.password", "12345678")
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", "3306")
    viper.SetDefault("database.dbname", "gober")
    viper.SetDefault("security.jwt.secret", "12345678")

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
