package initialize

import (
	"fmt"
	"log"
)

func Run() {
	// Tải cấu hình
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Không thể tải cấu hình: %v", err)
	}

	fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Database: user=%s, password=%s, host=%s, port=%d\n",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)

	r := InitRouter()

	r.Run() // 8080
}