test-handler: ## Test email handler
	go test -v ./internal/handler

test-handler-integration: ## Run integration tests (needs RESEND_API_KEY)
	go test -v -run Integration ./internal/handler

test-handler-coverage: ## Test with coverage report
	go test -coverprofile=coverage.out ./internal/handler
	go tool cover -html=coverage.out
```

---

## Expected Output
```
=== RUN   TestHandleMail_ValidPayload
--- PASS: TestHandleMail_ValidPayload (0.00s)
=== RUN   TestHandleMail_InvalidJSON
--- PASS: TestHandleMail_InvalidJSON (0.00s)
=== RUN   TestHandleMail_MissingRequiredFields
=== RUN   TestHandleMail_MissingRequiredFields/missing_to
--- PASS: TestHandleMail_MissingRequiredFields/missing_to (0.00s)
=== RUN   TestHandleMail_MissingRequiredFields/missing_from
--- PASS: TestHandleMail_MissingRequiredFields/missing_from (0.00s)
=== RUN   TestHandleMail_MissingRequiredFields/missing_subject
--- PASS: TestHandleMail_MissingRequiredFields/missing_subject (0.00s)
=== RUN   TestClassifyEmailError
=== RUN   TestClassifyEmailError/500_server_error
--- PASS: TestClassifyEmailError/500_server_error (0.00s)
=== RUN   TestClassifyEmailError/timeout_error
--- PASS: TestClassifyEmailError/timeout_error (0.00s)
=== RUN   TestHandleMail_Integration
--- SKIP: TestHandleMail_Integration (0.00s)
    handler_test.go:200: Skipping integration test: RESEND_API_KEY not set
PASS
coverage: 78.5% of statements
ok      yourproject/internal/handler    0.234s