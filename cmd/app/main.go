package main

import (
	config "Gober/configs"
	"Gober/internal/routers"
	"fmt"
	"log"
)

func main() {
  // Tải cấu hình
  cfg, err := config.LoadConfig()
  if err != nil {
		log.Fatalf("Không thể tải cấu hình: %v", err)
	}

  fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Database: user=%s, password=%s, host=%s, port=%d\n",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)
  
  r := routers.NewRouter()

  r.Run() // 8080
}

