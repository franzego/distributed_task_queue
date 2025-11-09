# Distributed Task Queue

A robust, scalable, and distributed task queue system built with Go. This project provides a simple yet powerful API for enqueuing jobs, which are then processed asynchronously by background workers. It's designed for reliability and performance, featuring API key authentication, rate limiting, and a clean separation between the web server and worker processes.

## Features
- **Asynchronous Job Processing**: Enqueue tasks via a REST API and let dedicated workers handle the processing in the background.
- **RESTful API**: A clean and simple API built with the Gin framework for managing jobs and administrative tasks.
- **Secure Endpoints**: Protects job submission endpoints with API Key authentication.
- **Rate Limiting**: Built-in token-bucket rate limiting on a per-API-key basis to prevent abuse.
- **Admin Controls**: Separate, token-protected endpoints for generating and managing API keys.
- **Persistent Storage**: Utilizes PostgreSQL to store job states, ensuring durability and data integrity.
- **Extensible Worker Logic**: Easily add new job types, with an initial implementation for sending emails via the Resend API.

## Technologies Used

| Technology | Description |
| :--- | :--- |
| [**Go**](https://golang.org/) | The core language used for its performance, concurrency, and simplicity. |
| [**Gin**](https://gin-gonic.com/) | A high-performance HTTP web framework for building the API server. |
| [**PostgreSQL**](https://www.postgresql.org/) | The relational database used for job persistence and API key storage. |
| [**pgx**](https://github.com/jackc/pgx) | High-performance PostgreSQL driver and toolkit for Go. |
| [**sqlc**](https://sqlc.dev/) | Generates type-safe Go code from SQL for database interactions. |
| [**Resend**](https://resend.com/) | Integrated for handling the "send_email" job type. |

## Getting Started

Follow these instructions to get the project up and running on your local machine.

### Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/franzego/distributed_task_queue.git
    cd distributed_task_queue
    ```

2.  **Install dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Set up the database**:
    Ensure you have a running PostgreSQL instance. You will need to create a database and run the schema migrations located in `db/schema.sql`.

4.  **Configure Environment Variables**:
    Create a `.env` file in the root of the project and add the necessary variables.

    ```sh
    # See the Environment Variables section for details
    cp .env.example .env
    ```

5.  **Run the application**:
    The system consists of two main components: the API server and the background worker. You need to run both in separate terminals.

    *   **Start the API Server**:
        ```bash
        go run ./cmd/server/main.go
        ```

    *   **Start the Worker**:
        ```bash
        go run ./cmd/worker/main.go
        ```

### Environment Variables

You need to create a `.env` file in the project root and populate it with the following variables.

| Variable | Purpose | Example |
| :--- | :--- | :--- |
| `DATABASE_URL` | The connection string for your PostgreSQL database. | `postgresql://user:password@localhost:5432/jobsdb?sslmode=disable` |
| `ADMIN_TOKEN` | A secret bearer token for accessing admin endpoints. | `my-super-secret-admin-token` |
| `RESEND_API_KEY`| Your API key from Resend for the email worker. | `re_123456789ABCDEF` |

## API Documentation

### Base URL
The API root path is hosted at the base URL where the server is running.
Example: `http://localhost:8080`

### Endpoints

#### Health Check
A public endpoint to verify that the API server is running.

#### `GET /health`
**Request**:
No payload required.

**Response**: `200 OK`
```json
{
  "message": "ok"
}
```

**Errors**:
- None

---

### Admin Endpoints
These endpoints are protected and require an `ADMIN_TOKEN`.

#### `POST /admin/api-keys`
Creates a new API key for accessing protected job endpoints.

**Request**:
- **Headers**: `Authorization: Bearer [ADMIN_TOKEN]`
- **Body**:
  ```json
  {
    "name": "My First App",
    "description": "API key for the primary application.",
    "prefix": "app"
  }
  ```

**Response**: `201 Created`
```json
{
  "id": "e6a5c1f0-a5c1-4b7e-8c1d-0f5a7d3b2a1c",
  "name": "My First App",
  "key": "app-a1b2c3d4e5f6...",
  "created_at": "2023-10-27T10:00:00Z",
  "warning": "Save this key securely. It wont be shown again"
}
```

**Errors**:
- `401 Unauthorized`: Missing or invalid `ADMIN_TOKEN`.
- `400 Bad Request`: Invalid request body.

---

#### `GET /admin/api-keys`
Lists all API keys created in the system.

**Request**:
- **Headers**: `Authorization: Bearer [ADMIN_TOKEN]`

**Response**: `200 OK`
```json
{
  "Message": "These are the keys",
  "Keys": [
    {
      "id": "e6a5c1f0-a5c1-4b7e-8c1d-0f5a7d3b2a1c",
      "name": "My First App",
      "key_hash": "hashed_key_value_1",
      "created_at": "2023-10-27T10:00:00Z",
      "last_used_at": "2023-10-27T10:05:00Z",
      "expires_at": null
    }
  ]
}
```

**Errors**:
- `401 Unauthorized`: Missing or invalid `ADMIN_TOKEN`.
- `500 Internal Server Error`: Failed to retrieve keys from the database.

---

### Job Endpoints
These endpoints are protected and require an `X-API-Key`.

#### `POST /jobs`
Enqueues a new job for asynchronous processing by a worker.

**Request**:
- **Headers**: `X-API-Key: [YOUR_API_KEY]`
- **Body**: The `payload` is a flexible JSON object that depends on the job `type`. For `send_email`:
  ```json
  {
    "type": "send_email",
    "payload": {
      "to": "recipient@example.com",
      "from": "sender@example.com",
      "subject": "Hello from the Task Queue!"
    }
  }
  ```

**Response**: `200 OK`
```json
{
  "message": "Message successfully received",
  "id": "Message of id: 1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed"
}
```

**Errors**:
- `401 Unauthorized`: Missing or invalid `X-API-Key`.
- `429 Too Many Requests`: Rate limit for the API key has been exceeded.
- `400 Bad Request`: Invalid or missing request body fields.
- `500 Internal Server Error`: Failed to enqueue the job.

---

#### `GET /jobs/{id}`
Retrieves the status and details of a specific job by its ID.
*Note: The route is `/jobs/:id`, but the current handler implementation incorrectly expects a query parameter `?id=...`. This documentation reflects the intended route parameter usage.*

**Request**:
- **Headers**: `X-API-Key: [YOUR_API_KEY]`
- **Path Parameter**: `id` (string, UUID)

**Response**: `200 OK`
```json
{
    "id": "1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed",
    "type": "send_email",
    "payload": {
        "to": "recipient@example.com",
        "from": "sender@example.com",
        "subject": "Hello from the Task Queue!"
    },
    "status": "completed",
    "attempts": 1,
    "max_attempts": 3,
    "created_at": "2023-10-27T12:00:00Z",
    "scheduled_at": "2023-10-27T12:00:00Z",
    "error_message": null
}
```

**Errors**:
- `401 Unauthorized`: Missing or invalid `X-API-Key`.
- `429 Too Many Requests`: Rate limit for the API key has been exceeded.
- `404 Not Found`: No job could be found with the provided ID.

## Contributing
Contributions are welcome! If you have suggestions for improvement or want to add new features, please feel free to open an issue or submit a pull request.

- 🍴 Fork the repository.
- ✨ Create a new branch (`git checkout -b feature/AmazingFeature`).
- 📝 Make your changes.
- ✅ Commit your changes (`git commit -m 'Add some AmazingFeature'`).
- 🚀 Push to the branch (`git push origin feature/AmazingFeature`).
- 📬 Open a pull request.

## License
This project is not licensed. Please contact the author for permissions.

## Author

**franzego**

- **LinkedIn**: [Your LinkedIn Profile](https://linkedin.com/in/your-username)
- **Twitter**: [@YourTwitterHandle](https://twitter.com/your-twitter-handle)

[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)