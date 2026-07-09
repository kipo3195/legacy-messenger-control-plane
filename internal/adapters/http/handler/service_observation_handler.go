package handler

import (
	"legacy-messenger-control-plane/internal/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceObservationHandler struct {
	serviceObservationUseCase usecase.ServiceObservationUsecase
}

func NewServiceObservationHandler(serviceObservationUseCase usecase.ServiceObservationUsecase) *ServiceObservationHandler {
	return &ServiceObservationHandler{
		serviceObservationUseCase: serviceObservationUseCase,
	}
}

func (h *ServiceObservationHandler) GetServiceListStatus(c *gin.Context) {
	result, err := h.serviceObservationUseCase.GetServiceList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ServiceObservationHandler) GetServiceStatus(c *gin.Context) {
	serviceName := c.Param("serviceName")

	result, err := h.serviceObservationUseCase.GetServiceStatus(c.Request.Context(), serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
