# Distributed Task Queue API

## Overview
This is a distributed task queue system built with Go. It features a RESTful API server using the Gin framework for enqueuing jobs and a separate background worker for processing them, utilizing PostgreSQL as the message broker and persistence layer.

## Features
- **Go**: High-performance, concurrent backend for the API and worker.
- **Gin**: A minimalistic and efficient web framework for routing and handling HTTP requests.
- **PostgreSQL**: Serves as a robust and persistent job queue database.
- **pgx/v5**: High-performance PostgreSQL driver and toolkit for Go.
- **sqlc**: Generates type-safe Go code from SQL queries for database interactions.
- **Background Worker**: Enables asynchronous and decoupled processing of submitted tasks.

## Getting Started
### Installation
1.  **Clone the repository**
    ```bash
    git clone https://github.com/franzego/distributed_task_queue.git
    cd distributed_task_queue
    ```

2.  **Install dependencies**
    ```bash
    go mod tidy
    ```

3.  **Set up the database**
    Ensure you have a running PostgreSQL instance. Create a database and execute the following schema. Note: a `schema.sql` file should be created based on your project's needs. Here is a sample based on the application's models:
    ```sql
    CREATE TABLE jobs (
        id UUID PRIMARY KEY,
        type VARCHAR(255) NOT NULL,
        payload JSONB NOT NULL,
        status VARCHAR(50) NOT NULL DEFAULT 'pending',
        attempts INT NOT NULL DEFAULT 0,
        max_attempts INT NOT NULL DEFAULT 3,
        error_message TEXT,
        scheduled_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

    CREATE INDEX idx_jobs_status_scheduled_at ON jobs (status, scheduled_at);
    ```

4.  **Configure Environment Variables**
    Update the `dburl` connection string in `cmd/server/main.go` and `cmd/worker/main.go` or modify the code to use environment variables.

5.  **Run the API Server**
    ```bash
    go run ./cmd/server/main.go
    ```
    The server will start on `http://localhost:8080` by default.

6.  **Run the Worker**
    In a separate terminal, start the worker process:
    ```bash
    go run ./cmd/worker/main.go
    ```

### Environment Variables
The application requires a database connection URL. It is recommended to modify the source to pull this from an environment variable instead of hardcoding it.

- `DATABASE_URL`: The connection string for your PostgreSQL database.
  - **Example**: `postgresql://user:password@localhost:5432/jobsdb?sslmode=disable`

## API Documentation
### Base URL
The API root path is `http://localhost:8080`.

### Endpoints
#### GET /health
Checks the health status of the API server.

**Request**:
No payload required.

**Response**:
```json
{
  "message": "ok"
}
```

**Errors**:
- None for this endpoint.

---
#### POST /jobs
Enqueues a new job for the worker to process.

**Request**:
The payload must be a JSON object with a `type` and a `payload`.
```json
{
  "type": "send_email",
  "payload": {
    "recipient": "test@example.com",
    "subject": "Hello from the Task Queue!",
    "body": "This is a test email."
  }
}
```

**Response**:
A success message including the unique ID of the enqueued job.
```json
{
  "message": "Message successfully received",
  "id": "Message of id: a1b2c3d4-e5f6-7890-1234-567890abcdef"
}
```

**Errors**:
- `400 Bad Request`: The request body is malformed, not valid JSON, or the `payload` field is empty.
- `500 Internal Server Error`: The server failed to create the job in the database.

---
#### GET /jobs/status
Retrieves the current status and details of a specific job by its ID.

**Request**:
The job ID must be provided as a URL query parameter.
- **Path**: `/jobs/status?id=<job-uuid>`

**Response**:
A JSON object containing the full details of the job.
```json
{
    "ID": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
    "Type": "send_email",
    "Payload": {
        "recipient": "test@example.com",
        "subject": "Hello from the Task Queue!",
        "body": "This is a test email."
    },
    "Status": "completed",
    "Attempts": 1,
    "MaxAttempts": 3,
    "ErrorMessage": null,
    "ScheduledAt": null,
    "CreatedAt": "2023-10-27T10:00:00Z",
    "UpdatedAt": "2023-10-27T10:00:05Z"
}
```

**Errors**:
- `404 Not Found`: The provided job ID does not exist in the database.