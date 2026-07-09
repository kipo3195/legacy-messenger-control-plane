package handler

import (
	"legacy-messenger-control-plane/internal/adapters/http/dto"
	"legacy-messenger-control-plane/internal/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceControlHandler struct {
	serviceControlUsecase usecase.ServiceControlUsecase
}

func NewServiceControlHandler(serviceControlUsecase usecase.ServiceControlUsecase) *ServiceControlHandler {
	return &ServiceControlHandler{
		serviceControlUsecase: serviceControlUsecase,
	}
}

func (h *ServiceControlHandler) ReDeploy(c *gin.Context) {

	serviceName := c.Param("serviceName")

	var req dto.ServiceRedeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	cmd := req.ToCommand(serviceName)

	result, err := h.serviceControlUsecase.ServiceRedeploy(c.Request.Context(), cmd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)

}
