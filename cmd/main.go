package main

import "Gober/internal/initialize"

func main() {
	// Chạy ứng dụng
	r, port := initialize.Run()

	r.Run(":" + port)
}

