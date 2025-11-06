package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	db "github.com/franzego/distributed_task_queue/db/sqlc"
	"github.com/franzego/distributed_task_queue/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	q Queue
}

func NewHandlerService(q Queue) *Handler {
	if q == nil {
		return nil
	}
	return &Handler{
		q: q,
	}
}

// Post Request To Enquque a Job
func (h *Handler) PostJob(c *gin.Context) {
	var req models.JobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Bad Request",
			Error:   err.Error(),
		})
		return
	}
	if len(req.Payload) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Payload cannot be empty",
		})
		return
	}
	uuid := uuid.New().String()
	job := db.Job{
		ID:          uuid,
		Type:        req.Type,
		Payload:     req.Payload,
		Status:      "pending",
		MaxAttempts: 3,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.q.Enqueue(ctx, job); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to Enqueue job",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessMessage{
		Message: "Message successfully received",
		ID:      fmt.Sprintf("Message of id: %s", job.ID),
	})

}

// Get Request To Get the Status of a particulat job using the uuid
func (h *Handler) GetStatus(c *gin.Context) {
	uuid := c.Query("id")
	job, err := h.q.GetJob(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "UUID could not be found",
		})
		return
	}
	// jsonJob, err := json.Marshal(&job)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, models.ErrorResponse{
	// 		Message: "Error in marshalling job",
	// 	})
	// 	return
	// }
	c.JSON(http.StatusOK, job)

}
