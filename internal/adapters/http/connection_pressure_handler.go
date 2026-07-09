package http

import (
	"legacy-messenger-control-plane/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConnectionPressureHandler struct {
	connectionPressureUsecase application.ConnectionPressureUsecase
}

func NewConnectionPressureHandler(connectionPressureUsecase application.ConnectionPressureUsecase) *ConnectionPressureHandler {
	return &ConnectionPressureHandler{
		connectionPressureUsecase: connectionPressureUsecase,
	}
}

func (h *ConnectionPressureHandler) GetConnectionStatus(c *gin.Context) {

	serviceName := c.Param("serviceName")
	result, err := h.connectionPressureUsecase.GetConnectionStatus(c.Request.Context(), serviceName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)

}
