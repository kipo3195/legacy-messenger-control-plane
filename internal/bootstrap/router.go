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
			// 관측
			services.GET("", handlers.ServiceObservation.GetServiceListStatus)
			services.GET("/:serviceName/status", handlers.ServiceObservation.GetServiceStatus)
			services.GET("/:serviceName/tasks", handlers.TaskObservation.GetTaskStatus)
			services.GET("/:serviceName/target-health", handlers.TargetHealth.GetTargetHealth)
			services.GET("/:serviceName/connection-pressure", handlers.ConnectionPressure.GetConnectionStatus)

			// 수집
			services.PUT("/:serviceName/tasks/:taskId/session-report", handlers.TaskSessionReport.PutTaskSessionReport)

			// 제어
			services.POST("/:serviceName/scale", handlers.ServiceScale.UpdateServiceDesiredCount)
			services.POST("/:serviceName/redeploy", handlers.ServiceControl.ReDeploy)

			// 판단
			services.POST("/:serviceName/scaling-evaluate", handlers.ServiceEvaluation.EvaluateScaling)

		}

	}

	return r
}
