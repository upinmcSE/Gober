package initialize

import (
	"log"
)

func Run() {

	// Khởi tạo cơ sở dữ liệu
	db := Init()
	if db == nil {
		log.Fatalf("Không thể khởi tạo cơ sở dữ liệu")
	}

	// Khởi tạo router
	r := InitRouter(db)

	r.Run() // 8080
}