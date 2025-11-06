package models

import "encoding/json"

type JobRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
type SuccessMessage struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}
