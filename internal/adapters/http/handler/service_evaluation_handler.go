package handler

import (
	"legacy-messenger-control-plane/internal/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceEvaluationHandler struct {
	serviceEvaluationUsecase usecase.ServiceEvaluationUsecase
}

func NewServiceEvaluationHandler(serviceEvaluationUsecase usecase.ServiceEvaluationUsecase) *ServiceEvaluationHandler {
	return &ServiceEvaluationHandler{
		serviceEvaluationUsecase: serviceEvaluationUsecase,
	}
}

func (h *ServiceEvaluationHandler) EvaluateScaling(c *gin.Context) {
	serviceName := c.Param("serviceName")

	result, err := h.serviceEvaluationUsecase.Evaluate(c.Request.Context(), serviceName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}
