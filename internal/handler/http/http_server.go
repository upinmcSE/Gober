package http

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/middleware"
	"Gober/pkg/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type HTTPServer interface {
	StartHTTPServer(wg *sync.WaitGroup)
}

type httpServer struct{}

func (h httpServer) StartHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// Đợi một chút để gRPC server khởi động trước
	time.Sleep(2 * time.Second)

	// grpc client
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	client := gober.NewGoberServiceClient(conn)

	r := gin.Default()

	v1 := r.Group("/v1/2025")

	httpLogger := logger.NewLoggerWithPath("../../internal/logs/http.log", "info")
	rateLimiterLogger := logger.NewLoggerWithPath("../../internal/logs/rate_limiter.log", "warning")
	v1.Use(
		middleware.ApikeyMiddleware(),
		middleware.CORSMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RateLimiterMiddleware(rateLimiterLogger),
	)

	v1.Use(gzip.Gzip(gzip.BestCompression))

	account := v1.Group("/accounts")
	accountHandler := NewAccountHandler(client)
	account.POST("/create", accountHandler.CreateHandler)
	account.POST("/session", accountHandler.CreateSessionHandler)
	account.GET("/:id", accountHandler.GetAccountHandler)

	//event := v1.Group("/events")

	//ticket := v1.Group("/tickets")

	log.Println("HTTP server listening at :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}

}

func NewHTTPServer() HTTPServer {
	return &httpServer{}
}
