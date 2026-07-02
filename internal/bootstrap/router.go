package bootstrap

import "github.com/gin-gonic/gin"

func NewRouter(handlers *Handlers) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/services/:serviceName/status", handlers.Service.GetStatus)
	}

	return r
}
