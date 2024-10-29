package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tiakavousi/codeplayground/pkg/container"
	"github.com/tiakavousi/codeplayground/pkg/executor"
)

const (
	defaultExecutionTimeout = 10 * time.Second
	defaultContainerImage   = "tayebe/repl"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO: Implement proper origin checking for production
			return true
		},
	}

	// Global variables for saved code functionality
	savedCodes = make(map[string]SavedCode)
	codesMutex sync.RWMutex

	// Global executor service
	execService *executor.Service
)

type SavedCode struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func main() {
	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize the executor service
	dockerImage := os.Getenv("DOCKER_IMAGE")
	if dockerImage == "" {
		dockerImage = defaultContainerImage
	}

	dockerRunner := container.NewDockerRunner(dockerImage)
	execService = executor.NewService(dockerRunner)

	// Initialize Gin router
	router := setupRouter()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"Upgrade",
		"Connection",
	}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	config.AllowWebSockets = true
	router.Use(cors.New(config))

	// Routes
	router.GET("/execute", handleWebSocket)
	router.POST("/save", handleSaveCode)
	router.GET("/share/:id", handleGetSavedCode)
	router.GET("/", handleHealthCheck)

	return router
}

func handleHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the backend!")
}

func handleWebSocket(c *gin.Context) {
	log.Printf(" handleWebSocket CALLED!!!")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Read initial request
	var req executor.ExecRequest
	if err := conn.ReadJSON(&req); err != nil {
		log.Println("JSON read error:", err)
		return
	}

	// Create execution context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), defaultExecutionTimeout)
	defer cancel()

	// Create channels for communication
	input := make(chan string, 5)
	output := make(chan string, 5)
	done := make(chan struct{})

	// Execute code
	go executeCode(ctx, req, input, output, done)

	// Handle WebSocket communication
	handleWebSocketCommunication(ctx, conn, input, output, done)
}

func executeCode(
	ctx context.Context, 
	req executor.ExecRequest, 
	input chan string, 
	output chan string, 
	done chan struct{}) {
	defer close(done)
	err := execService.ExecuteInteractive(ctx, req, input, output)
	if err != nil {
		select {
		case output <- "Execution error: " + err.Error():
		case <-ctx.Done():
		}
	}
}

func handleWebSocketCommunication(ctx context.Context, conn *websocket.Conn, input chan string, output chan string, done chan struct{}){
	// Handle input from WebSocket
	go func() {
		defer close(input)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
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

	// Handle output to WebSocket
	for {
		select {
		case line := <-output:
			if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		case <-done:
			// Clean shutdown
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		case <-ctx.Done():
			// Timeout
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Execution timeout"))
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
	log.Printf(" handleGetSavedCode CALLED!!! .GET: share/:id ")
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
