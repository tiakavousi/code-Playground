package main

import (
	"net/http"

	"github.com/tiakavousi/codeplayground/pkg/executor"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	// Manually attach Logger and Recovery if needed
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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

	// Run the web server
	router.Run(":8080")
}
