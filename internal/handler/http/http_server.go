package http

import (
	"Gober/internal/generated/grpc/gober"
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

	r.POST("/sessions", createSessionHandler(client))
	r.POST("/create", createHandler(client))
	r.GET("/:id", getAccountHandler(client))

	log.Println("HTTP server listening at :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}

}

func createSessionHandler(client gober.GoberServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req gober.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		resp, err := client.CreateSession(c, &req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

func createHandler(client gober.GoberServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req gober.CreateAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		resp, err := client.CreateAccount(c, &req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

func getAccountHandler(client gober.GoberServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Param("id")
		if accountID == "" {
			c.JSON(400, gin.H{"error": "Account ID is required"})
			return
		}

		// Convert accountID to int64
		id, err := strconv.ParseUint(accountID, 10, 64)

		resp, err := client.GetAccount(c, &gober.GetAccountRequest{AccountId: id})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

func NewHTTPServer() HTTPServer {
	return &httpServer{}
}
