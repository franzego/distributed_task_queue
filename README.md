# Distributed Task Queue API

## Overview
This project is a distributed task queue system built with Go. It features a RESTful API developed using the Gin framework for enqueuing and managing jobs, with PostgreSQL serving as the backend for persistent storage.

## Features
- **Go (Golang)**: Core backend language for high performance and concurrency.
- **Gin**: High-performance HTTP web framework for routing and handling API requests.
- **PostgreSQL**: Robust relational database used for job and API key persistence.
- **pgx/v5**: High-performance PostgreSQL driver for Go.
- **Background Worker**: A separate process for dequeuing and processing jobs asynchronously.

## Getting Started
### Installation
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/franzego/distributed_task_queue.git
    cd distributed_task_queue
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Set up the database:**
    Ensure you have a running PostgreSQL instance and create the necessary tables. Schema files are expected in a `db/schema.sql` directory (not provided in the context, but required for `sqlc`).

4.  **Create an environment file:**
    Create a `.env` file in the root directory and add the required environment variables.
    ```bash
    touch .env
    ```

5.  **Run the API server:**
    ```bash
    go run ./cmd/server/main.go
    ```

6.  **Run the background worker:**
    ```bash
    go run ./cmd/worker/main.go
    ```

### Environment Variables
All required environment variables should be placed in a `.env` file in the project root.

-   `DATABASE_URL`: The connection string for your PostgreSQL database.
-   `ADMIN_TOKEN`: A secret token for accessing administrative endpoints.

```env
DATABASE_URL="postgresql://user:password@localhost:5432/jobsdb?sslmode=disable"
ADMIN_TOKEN="your-super-secret-admin-token"
```

## API Documentation
### Base URL
The API is not versioned and all endpoints are relative to the base URL where the server is running.
Example: `http://localhost:8080`

### Endpoints
#### GET /health
A public endpoint to check the health status of the API server.

**Request**:
No payload required.

**Response**:
*Status: 200 OK*
```json
{
  "message": "ok"
}
```

**Errors**:
- None

---
#### POST /admin/api-keys
Creates a new API key. Requires admin authentication.

**Authentication**:
- Header: `Authorization: Bearer [ADMIN_TOKEN]`

**Request**:
```json
{
  "name": "WebApp Key",
  "description": "API key for the main web application",
  "prefix": "webapp"
}
```

**Response**:
*Status: 201 Created*
```json
{
  "id": "c6a2c2a0-4a6b-4e6e-9f3b-8f3a8e9e9e9e",
  "name": "WebApp Key",
  "created_at": "2023-10-27T10:00:00Z",
  "warning": "Save this key securely. It wont be shown again"
}
```

**Errors**:
-   `400 Bad Request`: Invalid JSON payload.
-   `401 Unauthorized`: Missing or invalid `ADMIN_TOKEN`.
-   `500 Internal Server Error`: Failed to generate or save the API key.

---
#### GET /admin/api-keys
Lists all existing API keys. Requires admin authentication.

**Authentication**:
- Header: `Authorization: Bearer [ADMIN_TOKEN]`

**Request**:
No payload required.

**Response**:
*Status: 200 OK*
```json
{
  "Message": "These are the keys",
  "Keys": [
    {
      "id": "c6a2c2a0-4a6b-4e6e-9f3b-8f3a8e9e9e9e",
      "name": "WebApp Key",
      "key_hash": "a1b2c3d4...",
      "created_at": "2023-10-27T10:00:00Z",
      "last_used": null,
      "expires_at": null
    }
  ]
}
```

**Errors**:
-   `401 Unauthorized`: Missing or invalid `ADMIN_TOKEN`.
-   `500 Internal Server Error`: Failed to retrieve API keys from the database.

---
#### POST /jobs
Enqueues a new job for processing by a background worker. Requires API key authentication.

**Authentication**:
- Header: `X-API-Key: [YOUR_API_KEY]`

**Request**:
The `payload` field can contain any valid JSON object.
```json
{
  "type": "send_email",
  "payload": {
    "recipient": "test@example.com",
    "subject": "Welcome!",
    "body": "Thank you for signing up."
  }
}
```

**Response**:
*Status: 200 OK*
```json
{
  "message": "Message successfully received",
  "id": "Message of id: a1b2c3d4-e5f6-a7b8-c9d0-e1f2a3b4c5d6"
}
```

**Errors**:
-   `400 Bad Request`: Invalid JSON payload or empty `payload`.
-   `401 Unauthorized`: Missing or invalid `X-API-Key`.
-   `500 Internal Server Error`: Failed to enqueue the job.

---
#### GET /jobs/:id
Retrieves the status and details of a specific job by its ID. Requires API key authentication.

**Authentication**:
- Header: `X-API-Key: [YOUR_API_KEY]`

**Request**:
The Job ID is passed as a URL parameter. Example: `/jobs/a1b2c3d4-e5f6-a7b8-c9d0-e1f2a3b4c5d6`

**Response**:
*Status: 200 OK*
```json
{
    "id": "a1b2c3d4-e5f6-a7b8-c9d0-e1f2a3b4c5d6",
    "type": "send_email",
    "payload": {
        "recipient": "test@example.com",
        "subject": "Welcome!",
        "body": "Thank you for signing up."
    },
    "status": "completed",
    "attempts": 1,
    "max_attempts": 3,
    "error_message": null,
    "created_at": "2023-10-27T11:00:00Z",
    "updated_at": "2023-10-27T11:00:05Z",
    "scheduled_at": null
}
```

**Errors**:
-   `401 Unauthorized`: Missing or invalid `X-API-Key`.
-   `404 Not Found`: No job could be found with the provided ID.

[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)