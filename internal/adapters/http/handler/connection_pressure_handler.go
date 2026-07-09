package handler

import (
	"legacy-messenger-control-plane/internal/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConnectionPressureHandler struct {
	connectionPressureUsecase usecase.ConnectionPressureUsecase
}

func NewConnectionPressureHandler(connectionPressureUsecase usecase.ConnectionPressureUsecase) *ConnectionPressureHandler {
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
