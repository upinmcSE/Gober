package initialize

import (
	config "Gober/configs"
	"log"
)

func Run() {
	// Khởi tạo cấu hình
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Không thể load cấu hình: %v", err)
	}

	// Khởi tạo cơ sở dữ liệu
	db := Init()
	if db == nil {
		log.Fatalf("Không thể khởi tạo cơ sở dữ liệu")
	}

	// Khởi tạo router
	r := InitRouter(db)

	r.Run() // 8080
}