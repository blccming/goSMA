package api

import (
	"net/http"
	"os"

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

func getMemoryMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().Memory)
}

func getNetworkMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().Network)
}

func getSystemMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().System)
}

func getTemperatureMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().Temperature)
}

func StartAPI() {
	r := gin.Default()

	r.GET("/health", healthCheck)
	r.GET("/metrics", getMetrics)
	r.GET("/metrics/cpu", getCpuMetrics)
	r.GET("/metrics/memory", getMemoryMetrics)
	r.GET("/metrics/network", getNetworkMetrics)
	r.GET("/metrics/system", getSystemMetrics)
	r.GET("/metrics/temperature", getTemperatureMetrics)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8485"
	}
	r.Run("0.0.0.0:" + port)
}
