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

func getNetworkMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().Network)
}

func getSystemMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().System)
}

func StartAPI() {
	r := gin.Default()

	r.GET("/health", healthCheck)
	r.GET("/metrics", getMetrics)
	r.GET("/metrics/cpu", getCpuMetrics)
	r.GET("/metrics/network", getNetworkMetrics)
	r.GET("/metrics/system", getSystemMetrics)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8485"
	}
	r.Run("0.0.0.0:" + port)
}
