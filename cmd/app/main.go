package main

import (
	"Gober/internal/routers"
)

func main() {
  r := routers.NewRouter()
  r.Run() // 8080
}

