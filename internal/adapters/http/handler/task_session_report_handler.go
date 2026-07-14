package handler

import (
	"legacy-messenger-control-plane/internal/adapters/http/dto"
	"legacy-messenger-control-plane/internal/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskSessionReportHandler struct {
	taskSessionReportUsecase usecase.TaskSessionReportUsecase
}

func NewTaskSessionReportHandler(taskSessionReportUsecase usecase.TaskSessionReportUsecase) *TaskSessionReportHandler {
	return &TaskSessionReportHandler{
		taskSessionReportUsecase: taskSessionReportUsecase,
	}
}

func (h *TaskSessionReportHandler) PutTaskSessionReport(c *gin.Context) {

	serviceName := c.Param("serviceName")
	taskID := c.Param("taskId")

	var req dto.TaskSessionReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	cmd := req.ToCommand(serviceName, taskID, *req.SessionCount)

	result, err := h.taskSessionReportUsecase.PutTaskSessionReport(c.Request.Context(), cmd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}
