package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"sync"
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

type SavedCode struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

var (
	savedCodes = make(map[string]SavedCode)
	codesMutex sync.RWMutex
)

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

	// New route for saving code
	router.POST("/save", handleSaveCode)

	// New route for retrieving saved code
	router.GET("/share/:id", handleGetSavedCode)

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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	input := make(chan string)
	output := make(chan string)
	done := make(chan struct{})

	// Goroutine for executing code
	go func() {
		defer close(done)
		err := executor.ExecuteInteractiveCode(ctx, req, input, output)
		if err != nil {
			select {
			case output <- "Execution error: " + err.Error():
			case <-ctx.Done():
			}
		}
	}()

	// Goroutine for reading from WebSocket
	go func() {
		defer cancel() // Ensure context is cancelled when this goroutine exits
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}
			select {
			case input <- string(message):
			case <-ctx.Done():
				return
			}
		}
	}()

	// Main loop for handling output
	for {
		select {
		case line := <-output:
			err := conn.WriteMessage(websocket.TextMessage, []byte(line))
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		case <-done:
			// Gracefully close the WebSocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("Error during closing WebSocket: %v", err)
			}
			return
		case <-ctx.Done():
			// Gracefully close the WebSocket connection on timeout
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Execution timeout"))
			if err != nil {
				log.Printf("Error during closing WebSocket on timeout: %v", err)
			}
			return
		}
	}
}

func handleSaveCode(c *gin.Context) {
	var req SavedCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := generateUniqueID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique ID"})
		return
	}

	codesMutex.Lock()
	savedCodes[id] = req
	codesMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func handleGetSavedCode(c *gin.Context) {
	id := c.Param("id")

	codesMutex.RLock()
	savedCode, exists := savedCodes[id]
	codesMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code not found"})
		return
	}

	c.JSON(http.StatusOK, savedCode)
}

func generateUniqueID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
