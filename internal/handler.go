package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/franzego/distributed_task_queue/authutil"
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

// Post Request For Admin to create Api Keys
func (h *Handler) PostAdminApiKey(c *gin.Context) {
	var req struct {
		Name        string `json:"name"` //name of the api key
		Description string `json:"description"`
		Prefix      string `json:"prefix"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid Request",
			Error:   err.Error(),
		})
		return
	}
	key, err := authutil.KeyGenerator(req.Prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to generate Api key",
		})
		return
	}
	hashKey := authutil.HashApiKeys(key)

	newKey, err := h.q.CreateAPIKey(c.Request.Context(), db.CreateAPIKeyParams{
		ID:      uuid.New().String(),
		Name:    req.Name,
		KeyHash: hashKey,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to create Api key",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":         newKey.ID,
		"name":       newKey.Name,
		"key":        newKey,
		"created_at": newKey.CreatedAt,
		"warning":    "Save this key securely. It wont be shown again",
	})
}

// Get Request For Admin to list all Api keys
func (h *Handler) GetApiKeys(c *gin.Context) {
	keys, err := h.q.ListAPIKeys(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Could not list Api Keys",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "These are the keys",
		"Keys":    keys,
	})

}
