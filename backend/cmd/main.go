package main

import (
	"net/http"
	"time"

	"github.com/tiakavousi/codeplayground/pkg/executor"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	router := gin.New()

	// Attach Logger and Recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Enable CORS with custom configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow your frontend's origin
		AllowMethods:     []string{"POST", "OPTIONS"},       // Allow POST and OPTIONS requests
		AllowHeaders:     []string{"Origin", "Content-Type"},// Allow necessary headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,                             // Allow credentials if needed
		MaxAge:           12 * time.Hour,
	}))

	// Handle OPTIONS preflight request for /execute
	router.OPTIONS("/execute", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Status(http.StatusOK)
	})

	// Route to handle code execution requests
	router.POST("/execute", func(c *gin.Context) {
		var req executor.ExecRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call the execution function
		output, err := executor.ExecuteCode(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"output": output})
	})

	// Basic route to test if server is working
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the backend!")
	})

	// Run the web server on port 8080
	router.Run(":8080")
}
