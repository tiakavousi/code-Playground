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
		return true // Be careful with this in production
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
		log.Println(err)
		return
	}
	defer conn.Close()

	// Read the initial request
	var req executor.ExecRequest
	err = conn.ReadJSON(&req)
	if err != nil {
		log.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := make(chan string)
	output := make(chan string)

	go func() {
		err := executor.ExecuteInteractiveCode(ctx, req, input, output)
		if err != nil {
			log.Println(err)
			conn.WriteMessage(websocket.TextMessage, []byte("Execution error: "+err.Error()))
		}
		close(output)
	}()

	// Handle incoming messages (user input)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			input <- string(message)
		}
	}()

	// Send output back to the client
	for line := range output {
		err := conn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println(err)
			return
		}
	}
}
