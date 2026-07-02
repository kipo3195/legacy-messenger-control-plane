package http

import (
	"net/http"
	"legacy-messenger-control-plane/internal/ports"
)

type ServiceHandler struct { 
	statusUseCase ports.ServiceStatusUseCase
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