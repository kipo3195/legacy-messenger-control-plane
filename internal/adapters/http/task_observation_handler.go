package http

import (
	"legacy-messenger-control-plane/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskObservationHandler struct {
	taskObservationUsecase application.TaskObservationUsecase
}

func NewTaskObservationHandler(taskObservationUsecase application.TaskObservationUsecase) *TaskObservationHandler {
	return &TaskObservationHandler{
		taskObservationUsecase: taskObservationUsecase,
	}
}

func (h *TaskObservationHandler) GetTaskStatus(c *gin.Context) {

	serviceName := c.Param("serviceName")
	desiredStatus := c.Param("desiredStatus")
	result, err := h.taskObservationUsecase.GetTaskStatus(c.Request.Context(), serviceName, desiredStatus)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}
