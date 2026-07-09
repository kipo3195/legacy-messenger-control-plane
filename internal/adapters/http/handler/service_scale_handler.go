package handler

import (
	"legacy-messenger-control-plane/internal/adapters/http/dto"
	"legacy-messenger-control-plane/internal/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceScaleHandler struct {
	serviceScaleUsecase usecase.ServiceScaleUsecase
}

func NewServiceScaleHandler(serviceScaleUsecase usecase.ServiceScaleUsecase) *ServiceScaleHandler {
	return &ServiceScaleHandler{
		serviceScaleUsecase: serviceScaleUsecase,
	}
}

func (h *ServiceScaleHandler) UpdateServiceDesiredCount(c *gin.Context) {

	serviceName := c.Param("serviceName")

	var req dto.ServiceScaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	cmd := req.ToCommand(serviceName)

	result, err := h.serviceScaleUsecase.UpdateServiceDesirdCount(c.Request.Context(), cmd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}
