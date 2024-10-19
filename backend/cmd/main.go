package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tiakavousi/codeplayground/pkg/executor"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	router := gin.New()

	// Attach Logger and Recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Enable CORS with custom configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "Upgrade", "Connection"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	config.AllowWebSockets = true
	router.Use(cors.New(config))

	// WebSocket route for interactive code execution
	router.GET("/execute", func(c *gin.Context) {
		handleWebSocket(c.Writer, c.Request)
	})

	// Keeping the old POST route for non-interactive execution
	router.POST("/execute", func(c *gin.Context) {
		var req executor.ExecRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	var req executor.ExecRequest
	err = conn.ReadJSON(&req)
	if err != nil {
		log.Println("JSON read error:", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	input := make(chan string)
	output := make(chan string)
	done := make(chan struct{})
	writeChannel := make(chan string)

	// Goroutine for writing to WebSocket
	go func() {
		for {
			select {
			case message := <-writeChannel:
				err := conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Println("WebSocket write error:", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Goroutine for executing code
	go func() {
		defer close(done)
		err := executor.ExecuteInteractiveCode(ctx, req, input, output)
		if err != nil {
			writeChannel <- "Execution error: " + err.Error()
		}
	}()

	// Goroutine for reading from WebSocket
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				cancel() // Cancel the context to stop execution
				return
			}
			select {
			case input <- string(message):
			case <-done:
				return
			}
		}
	}()

	// Main loop for handling output
	for {
		select {
		case line := <-output:
			select {
			case writeChannel <- line:
			case <-done:
				return
			}
		case <-done:
			return
		}
	}
}
