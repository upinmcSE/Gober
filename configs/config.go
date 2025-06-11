package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// Config chứa các giá trị cấu hình
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Security SecurityConfig
}

// ServerConfig chứa cấu hình server
type ServerConfig struct {
    Host string
    Port int
}

type SecurityConfig struct {
    SecretKey string
}

// DatabaseConfig chứa cấu hình database
type DatabaseConfig struct {
    User     string
    Password string
    Host     string
    Port     int
    DbName   string
}

// Biến toàn cục để lưu cấu hình
var (
    config     Config
    configOnce sync.Once
)

// LoadConfig khởi tạo và đọc cấu hình từ Viper
func LoadConfig() error {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")

    // Thiết lập giá trị mặc định
    viper.SetDefault("server.host", "127.0.0.1")
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("database.user", "12345678")
    viper.SetDefault("database.password", "12345678")
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", 3306)
    viper.SetDefault("database.dbname", "gober")
    viper.SetDefault("security.jwt.secret", "12345678")

    // Đọc tệp cấu hình
    if err := viper.ReadInConfig(); err != nil {
        return fmt.Errorf("lỗi khi đọc tệp cấu hình: %w", err)
    }

    // Gán giá trị vào struct Config
    configOnce.Do(func() {
        config.Server.Host = viper.GetString("server.host")
        config.Server.Port = viper.GetInt("server.port")
        config.Database.User = viper.GetString("database.user")
        config.Database.Password = viper.GetString("database.password")
        config.Database.Host = viper.GetString("database.host")
        config.Database.Port = viper.GetInt("database.port")
        config.Database.DbName = viper.GetString("database.dbname")
        config.Security.SecretKey = viper.GetString("security.jwt.secret")
    })

    return nil
}

func GetConfig() Config {
    return config
}