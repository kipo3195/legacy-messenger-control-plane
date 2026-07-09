package http

import (
	"legacy-messenger-control-plane/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TargetHealthHandler struct {
	targetHealthUsecase application.TargetHealthUsecase
}

func NewTargetHealthHandler(targetHealthUsecase application.TargetHealthUsecase) *TargetHealthHandler {
	return &TargetHealthHandler{
		targetHealthUsecase: targetHealthUsecase,
	}
}

func (h *TargetHealthHandler) GetTargetHealth(c *gin.Context) {

	serviceName := c.Param("serviceName")
	result, err := h.targetHealthUsecase.GetTargetHealth(c.Request.Context(), serviceName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)

}
