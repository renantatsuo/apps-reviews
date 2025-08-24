# Apps Reviews

The system continuously polls the Apple App Store for new reviews of configured apps and provides a web interface to search and view them.

## Project Structure

More information can be found on [web/README.md](web/README.md) and [server/README.md](server/README.md)

## Features

- **Continuous Polling**: Automatically fetches new reviews at configurable intervals (default: 30 seconds)
- **Time-based Filtering**: Only stores and displays reviews from the last 48 hours (configurable)
- **Multi-app Support** \*: Monitor multiple Apple App Store apps simultaneously (for now only accepts app ids from service configuration)
- **REST API**: HTTP API for accessing review data
- **Web App**: React UI for searching and viewing reviews
- **Persistent Storage**: Local file-based storage for review data

## Quick Start

### Prerequisites

- Go 1.24.2+
- Node.js 18+
- npm

### Installation

```bash
make init
```

### Running

```bash
# Start both server and web interface
make dev

# Or start them separately:
make dev-server  # Go server on :8080
make dev-web     # React dev server on :5173
```

### Example Usage

```bash
# Configure app IDs to monitor
export APP_IDS="389801252,1234567890"

# Set polling interval to 1 minute
export POLLING_INTERVAL="1m"

# Start the server
make dev-server
```

Then visit `http://localhost:5173` to search for reviews by App ID.

## API Reference

### GET /api/reviews/{appID}

Fetch reviews for a specific Apple App Store app ID.

**Response:**

```json
{
  "data": [
    {
      "id": "review-id",
      "author": {
        "name": "User Name",
        "uri": "user-profile-uri"
      },
      "title": "Review Title",
      "content": "Review content...",
      "rating": 5,
      "updated": "2024-01-01T12:00:00Z"
    }
  ]
}
```

**Status Codes:**

- `200` - Success
- `404` - App not found or no reviews available
- `500` - Internal server error

## TODO

### Missing Features & Improvements

- [ ] **Concurrency Safety** - for now the polling engine do not support multiple instances
- [ ] **Testing**
- [ ] **API**

  - [ ] Pagination support for larger apps
  - [ ] Live fetch for new app IDs (on-demand add new apps)
