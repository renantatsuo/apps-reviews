# app-review server

A built with go HTTP server to provide an API for retrieving Apple App Store reviews.

## Architecture

The server follows a clean, layered architecture with clear separation of concerns:

**Application Layer (`cmd/server/`)**

- Entry point with dependency injection and graceful shutdown handling
- Configures logging, environment variables, and starts background services

**HTTP Layer (`internal/server/`)**

- RESTful API server with route handling
- CORS middleware for cross-origin requests
- Request/response transformation and error handling

**Business Logic Layer (`internal/reviews/`)**

- Reviews client that orchestrates Apple API integration
- Background polling service that fetches new reviews periodically
- Data persistence layer that manages file storage and in-memory caching

**Data Layer (`internal/models/`, `internal/store/`)**

- Domain models with transformation logic between Apple API and internal formats
- Thread-safe in-memory store for fast access
- File-based persistence for data durability

**External Integration (`pkg/apple/`)**

- Apple App Store RSS feed client
- Data structures and parsing logic for Apple's review format

## Running the Server

Run the server directly with Go:

```bash
go run ./cmd/server/...
```

The server will start on port 8080 by default and begin polling Apple's RSS feeds for the configured app IDs.

## API Endpoints

### Get Reviews

```
GET /reviews/{appID}
```

Returns reviews for the specified Apple App ID.

**Response:**

```json
{
  "data": [
    {
      "id": "review-id",
      "author": {
        "name": "Author Name",
        "uri": "author-uri"
      },
      "title": "Review Title",
      "content": "Review content...",
      "rating": 5,
      "updated": "2024-01-01T12:00:00Z"
    }
  ]
}
```

## Features

- **Background Polling**: Automatically fetches new reviews from Apple's RSS feeds at configured intervals
- **Data Persistence**: Reviews are cached both in-memory for fast access and on disk for durability
- **Graceful Shutdown**: Server handles SIGINT/SIGTERM signals for clean shutdown
- **CORS Support**: Cross-origin requests are supported for web applications

## Configuration

The server can be configured using environment variables:

### Environment Variables

| Variable             | Description                                      | Default                | Example                       |
| -------------------- | ------------------------------------------------ | ---------------------- | ----------------------------- |
| `PORT`               | The port number on which the server will listen  | `8080`                 | `PORT=3000`                   |
| `LOG_LEVEL`          | The logging level for the application            | `debug`                | `LOG_LEVEL=info`              |
| `APP_IDS`            | Comma-separated list of Apple App IDs to monitor | `1458862350,389801252` | `APP_IDS=123456789,987654321` |
| `STORE_DIR`          | Directory path for storing cached review data    | `data`                 | `STORE_DIR=/tmp/reviews`      |
| `REVIEWS_TIME_LIMIT` | How far back to fetch reviews from Apple         | `48h`                  | `REVIEWS_TIME_LIMIT=72h`      |
| `POLLING_INTERVAL`   | How often to poll Apple for new reviews          | `30s`                  | `POLLING_INTERVAL=5m`         |

#### Log Levels

Available log levels:

| Level   | Description                  | Use Case                                   |
| ------- | ---------------------------- | ------------------------------------------ |
| `debug` | Detailed debug information   | Development and troubleshooting            |
| `info`  | General information messages | Normal operation monitoring                |
| `warn`  | Warning messages             | Potential issues that don't stop operation |
| `error` | Error messages               | Errors that affect functionality           |
