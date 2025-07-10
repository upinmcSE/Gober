package http

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/middleware"
	"Gober/pkg/logger"
	"Gober/utils/response"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"strconv"
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

	account := v1.Group("/accounts")
	{
		account.POST("/login", createSessionHandler(client))
		account.POST("/register", createHandler(client))
		account.GET("/:id", middleware.AuthMiddleware(), getAccountHandler(client))
	}

	log.Println("HTTP server listening at :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}

}

func createSessionHandler(client gober.GoberServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req gober.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ResponseError(c, 400, &response.AppError{
				Message: "Invalid request",
				Code:    response.ErrCodeBadRequest,
			})
			return
		}

		resp, err := client.CreateSession(c, &req)
		if err != nil {
			response.HandleGrpcError(c, err)
			return
		}

		response.ResponseSuccess(c, 200, "Session created successfully", resp)
	}
}

func createHandler(client gober.GoberServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req gober.CreateAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ResponseError(c, 400, &response.AppError{
				Message: "Invalid request",
				Code:    response.ErrCodeBadRequest,
			})
			return
		}

		resp, err := client.CreateAccount(c, &req)
		if err != nil {
			response.HandleGrpcError(c, err)
			return
		}

		response.ResponseSuccess(c, 201, "Account created successfully", resp)
	}
}

func getAccountHandler(client gober.GoberServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Param("id")
		if accountID == "" {
			response.ResponseError(c, 400, &response.AppError{
				Message: "Account ID is required",
				Code:    response.ErrCodeBadRequest,
			})
			return
		}

		// Convert accountID to int64
		id, err := strconv.ParseUint(accountID, 10, 64)

		resp, err := client.GetAccount(c, &gober.GetAccountRequest{AccountId: id})
		if err != nil {
			response.HandleGrpcError(c, err)
			return
		}

		response.ResponseSuccess(c, 200, "Account retrieved", resp)
	}
}

func NewHTTPServer() HTTPServer {
	return &httpServer{}
}
