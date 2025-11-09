package handler

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestHandleEmail_ValidPayload(t *testing.T) {
	handler := &EmailHandler{
		ApiKey: "test_key",
	}
	payload := EmailPayload{
		To:      "test@example.com",
		From:    "sender@example.com",
		Subject: "Test Subject",
	}
	ctx := context.Background()
	payloadJson, _ := json.Marshal(payload)
	err := handler.HandleMail(ctx, payloadJson)
	if err != nil && err.Error() == "there was a problem with the payload request" {
		t.Fatal("Payload parsing failed")
	}
}
func TestHandleEmail_InvalidJson(t *testing.T) {
	handler := &EmailHandler{
		ApiKey: "test_key",
	}
	payload := EmailPayload{
		To:      "test@example.com",
		From:    "",
		Subject: "Test Subject",
	}
	ctx := context.Background()
	payloadJson, _ := json.Marshal(payload)
	err := handler.HandleMail(ctx, payloadJson)
	if err == nil {
		t.Fatal("expected error for invalid json")
	}
	if !strings.Contains(err.Error(), "missing requrired fields") {
		t.Fatalf("missing requrired fields error, got: %v", err)
	}
}
func TestHandleEmail_MissingFields(t *testing.T) {
	tests := []struct {
		name    string
		payload EmailPayload
	}{
		{
			name: "missing from",
			payload: EmailPayload{
				To:      "musa",
				Subject: "happy",
			},
		},
		{
			name: "missing to",
			payload: EmailPayload{
				From:    "musa",
				Subject: "happy",
			},
		},
		{
			name: "missing subject",
			payload: EmailPayload{
				From: "musa",
				To:   "happy",
			},
		},
		{
			name:    "missing payload",
			payload: EmailPayload{},
		},
	}
	handler := &EmailHandler{ApiKey: "test_key"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadJson, _ := json.Marshal(tt.payload)
			err := handler.HandleMail(context.Background(), payloadJson)
			if err == nil {
				t.Fatal("Expected error for missing fiels and got nil")
			}
			if !strings.Contains(err.Error(), "missing requrired fields") {
				t.Fatalf("Expected missing fields error, got: %v", err)
			}
		})
	}
}
func TestNewEmailHandlerService(t *testing.T) {
	os.Setenv("RESEND_API_KEY", "123456789")
	defer os.Unsetenv("RESEND_API_KEY")
	handler := NewEmailHandlerService()
	if handler == nil {
		t.Fatal("Handler should not be nil")
	}
	if handler.ApiKey != "123456789" {
		t.Fatalf("Expected 123456789, but got these %s", handler.ApiKey)
	}
}
func TestMissingEnvVariable(t *testing.T) {
	os.Unsetenv("RESEND_API_KEY")
	handler := NewEmailHandlerService()
	if handler == nil {
		t.Fatal("Handler should not be nil")
	}
	if handler.ApiKey != "" {
		t.Fatalf("api key should be empty but instead have this- %s", handler.ApiKey)
	}
}
func TestHandleContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	payload := EmailPayload{
		To:      "test@example.com",
		From:    "sender@example.com",
		Subject: "Test",
	}
	payloadJSON, _ := json.Marshal(payload)
	handler := NewEmailHandlerService()

	err := handler.HandleMail(ctx, payloadJSON)
	if err == nil {
		t.Fatal("Expecting a context cancelled error")
	}

}
func TestHandleMail_Integration(t *testing.T) {
	// Skip if no API key (for CI/CD)
	err := godotenv.Load()
	if err != nil {
		t.Skip("Skipping integration as we couldnt load env variable")
	}
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" || apiKey == "test_key" {
		t.Skipf("Skipping integration test: could only see this: %s", apiKey)
	}

	handler := NewEmailHandlerService()

	payload := EmailPayload{
		To:      "delivered@resend.dev", // Resend's test address
		From:    "onboarding@resend.dev",
		Subject: "Integration Test Email",
	}
	payloadJSON, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = handler.HandleMail(ctx, payloadJSON)

	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	t.Log("✅ Email sent successfully in integration test")
}
