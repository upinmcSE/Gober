package initialize

import (
	config "Gober/configs"
	"Gober/database"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() *gorm.DB {

	// Tải cấu hình từ tệp config.yaml
	cfg := config.GetConfig()

	// Khởi tạo kết nối đến cơ sở dữ liệu
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DbName)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Lỗi khi kết nối đến cơ sở dữ liệu: %v\n", err)
		return nil
	}

	// Tự động tạo bảng nếu chưa tồn tại
	if err := database.DBMigrator(db); err != nil {
		log.Fatalf("Không thể migrate db: %v", err)
	}
	return db
}