package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(handlers *Handlers) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	v1 := r.Group("/api/v1")
	{
		services := v1.Group("/services")
		{
			services.GET("", handlers.ServiceObservation.GetServiceListStatus)
			services.GET("/:serviceName/status", handlers.ServiceObservation.GetServiceStatus)
			services.GET("/:serviceName/task", handlers.TaskObservation.GetTaskStatus)

			services.GET("/:serviceName/target-health", handlers.TargetHealth.GetTargetHealth)
			services.GET("/:serviceName/connection-pressure", handlers.ConnectionPressure.GetConnectionStatus)

		}
	}

	return r
}
