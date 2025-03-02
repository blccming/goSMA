package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getCpuMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().CPU)
}

func getSystemMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getEndpointData().System)
}

func StartAPI() {
	router := gin.Default()
	router.GET("/metrics/cpu", getCpuMetrics)
	router.GET("/metrics/system", getSystemMetrics)

	router.Run("localhost:8080")
}
