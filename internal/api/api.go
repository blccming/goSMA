package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "API service is running",
	})
}

func getMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData())
}

func getCpuMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().CPU)
}

func getSystemMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().System)
}

func StartAPI() {
	r := gin.Default()
	r.GET("/health", healthCheck)

	r.GET("/metrics", getMetrics)
	r.GET("/metrics/cpu", getCpuMetrics)
	r.GET("/metrics/system", getSystemMetrics)

	r.Run("localhost:8080")
}
