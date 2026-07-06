package http

import (
	"legacy-messenger-control-plane/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	statusUseCase application.ServiceStatusUsecase
}

func (h *ServiceHandler) GetServiceStatus(c *gin.Context) {
	serviceName := c.Param("serviceName")

	result, err := h.statusUseCase.Execute(c.Request.Context(), serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
