package initialize

import (
	"fmt"

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

// LoadConfig khởi tạo và đọc cấu hình từ Viper
func LoadConfig() (Config, error) {
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
		return Config{}, fmt.Errorf("lỗi khi đọc tệp cấu hình: %w", err)
	}

	// Gán giá trị vào struct Config
	var cfg Config
	cfg.Server.Host = viper.GetString("server.host")
	cfg.Server.Port = viper.GetInt("server.port")
	cfg.Database.User = viper.GetString("database.user")
	cfg.Database.Password = viper.GetString("database.password")
	cfg.Database.Host = viper.GetString("database.host")
	cfg.Database.Port = viper.GetInt("database.port")
	cfg.Database.DbName = viper.GetString("database.dbname")
	cfg.Security.SecretKey = viper.GetString("security.jwt.secret")

	return cfg, nil
}