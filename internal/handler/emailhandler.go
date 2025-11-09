package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/resend/resend-go/v3"
)

var ErrInvalidPayload = errors.New("invalid email payload")

type EmailPayload struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
}

type EmailHandler struct {
	ApiKey string //resend api key
	// httpclient *http.Client
}

func NewEmailHandlerService() *EmailHandler {
	return &EmailHandler{
		ApiKey: os.Getenv("RESEND_API_KEY"),
		// httpclient: &http.Client{
		// 	Timeout: 10 * time.Second,
		// },
	}
}
func (e *EmailHandler) HandleMail(ctx context.Context, payload json.RawMessage) error {
	var req EmailPayload
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return fmt.Errorf("there was a problem with the payload request: %v", err)
	}
	if req.From == "" || req.To == "" || req.Subject == "" {
		return fmt.Errorf("missing requrired fields: %v", ErrInvalidPayload)
	}
	if err := e.Sendemail(ctx, req); err != nil {
		return err
	}
	return nil
}

func (e *EmailHandler) Sendemail(ctx context.Context, mail EmailPayload) error {
	client := resend.NewClient(e.ApiKey)

	params := &resend.SendEmailRequest{
		From:    mail.From,
		To:      []string{mail.To},
		Html:    "<strong>hello To you all</strong>",
		Subject: mail.Subject,
		// Cc:      []string{"cc@example.com"},
		// Bcc:     []string{"bcc@example.com"},
		// ReplyTo: "replyto@example.com",
	}
	responseEmail, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return ClassifyEmailError(err)
	}
	log.Printf("email was sent successfully: id=%s", responseEmail.Id)
	return nil

}

type RetriableError struct {
	Msg string
}

func (e *RetriableError) Error() string {
	return e.Msg
}

func ClassifyEmailError(err error) error {
	msg := err.Error()

	if strings.Contains(msg, "500") ||
		strings.Contains(msg, "502") ||
		strings.Contains(msg, "503") ||
		strings.Contains(msg, "timeout") {
		return &RetriableError{
			Msg: fmt.Sprintf("server error: %v", err),
		}
	}
	return &RetriableError{
		Msg: fmt.Sprintf("unknown error: %v", err),
	}
}
