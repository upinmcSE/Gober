package initialize

import (
	"Gober/configs"
	"log"

	"github.com/gin-gonic/gin"
)

func Run() (*gin.Engine, string) {
	// 1> Read config -> environment variables
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	// 2> Initialize database connection
	db, err := InitDB(&config)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	// 3> Initialize router
	r := InitRouter(db, &config)

	// 4> Initialize other services if needed (e.g., cache, message queue, etc.)
	
	return r, config.Server.Port
}