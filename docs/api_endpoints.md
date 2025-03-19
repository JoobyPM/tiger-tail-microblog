# Tiger-Tail API Endpoints

This document describes the API endpoints available in the Tiger-Tail Microblog application, following our [Tiger Style](tiger_style.md) principles of safety, performance, and developer experience.

## Overview

Tiger-Tail exposes a RESTful API with the following characteristics:

- All endpoints return JSON responses
- Authentication is done via Basic Auth
- All timestamps are in ISO 8601 format (UTC)
- Pagination is supported for list endpoints
- Rate limiting is applied to prevent abuse
- Error responses follow a consistent format

## Base URL

For local development: `http://localhost:8080`
For production: `https://your-deployment-url`

## Health Check Endpoints

### GET /livez

Liveness probe for Kubernetes.

**Response (200 OK):**
```json
{
  "status": "alive"
}
```

### GET /readyz

Readiness probe for Kubernetes.

**Response (200 OK):**
```json
{
  "status": "ready",
  "dependencies": {
    "database": "connected",
    "redis": "connected"
  }
}
```

## Posts Endpoints

### GET /posts

Returns a list of posts, with optional pagination.

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of posts per page (default: 10)

**Response (200 OK):**
```json
{
  "posts": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "content": "This is a post about Tiger-Tail!",
      "created_at": "2025-03-18T12:00:00Z",
      "updated_at": "2025-03-18T12:00:00Z"
    },
    {
      "id": "223e4567-e89b-12d3-a456-426614174001",
      "content": "Another interesting post.",
      "created_at": "2025-03-18T11:30:00Z",
      "updated_at": "2025-03-18T11:30:00Z"
    }
  ],
  "pagination": {
    "total": 42,
    "page": 1,
    "limit": 10,
    "pages": 5
  }
}
```

### GET /posts/{id}

Returns a specific post by ID.

**Path Parameters:**
- `id`: Post ID (UUID)

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "content": "This is a post about Tiger-Tail!",
  "created_at": "2025-03-18T12:00:00Z",
  "updated_at": "2025-03-18T12:00:00Z"
}
```

**Response (404 Not Found):**
```json
{
  "error": "not_found",
  "message": "Post not found"
}
```

### POST /posts

Creates a new post. Requires authentication.

**Request Headers:**
- `Authorization`: Basic Auth header

**Request Body:**
```json
{
  "content": "This is a new post about Tiger-Tail!"
}
```

**Response (201 Created):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "content": "This is a new post about Tiger-Tail!",
  "created_at": "2025-03-18T12:00:00Z",
  "updated_at": "2025-03-18T12:00:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "validation_error",
  "message": "Content is required",
  "details": {
    "content": "required"
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "unauthorized",
  "message": "Authentication required"
}
```

### PUT /posts/{id}

Updates an existing post. Requires authentication.

**Path Parameters:**
- `id`: Post ID (UUID)

**Request Headers:**
- `Authorization`: Basic Auth header

**Request Body:**
```json
{
  "content": "Updated post content"
}
```

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "content": "Updated post content",
  "created_at": "2025-03-18T12:00:00Z",
  "updated_at": "2025-03-18T12:05:00Z"
}
```

### DELETE /posts/{id}

Deletes a post. Requires authentication.

**Path Parameters:**
- `id`: Post ID (UUID)

**Request Headers:**
- `Authorization`: Basic Auth header

**Response (204 No Content)**

## Error Handling

All API endpoints follow a consistent error response format:

```json
{
  "error": "error_code",
  "message": "Human-readable error message",
  "details": {
    // Optional field-specific error details
  }
}
```

### Common Error Codes

| Status Code | Error Code       | Description                        |
|-------------|------------------|------------------------------------|
| 400         | validation_error | Invalid request parameters or body |
| 401         | unauthorized     | Authentication required            |
| 403         | forbidden        | Insufficient permissions           |
| 404         | not_found        | Resource not found                 |
| 429         | rate_limited     | Too many requests                  |
| 500         | internal_error   | Server error                       |

## Rate Limiting

To ensure system stability and prevent abuse, the API implements rate limiting:

- 100 requests per minute per IP address
- 1000 requests per hour per IP address

When rate limited, the API returns a 429 Too Many Requests response with headers:

- `X-RateLimit-Limit`: The rate limit ceiling
- `X-RateLimit-Remaining`: The number of requests left for the time window
- `X-RateLimit-Reset`: The remaining window before the rate limit resets (in seconds)

## Pagination

List endpoints support pagination with the following query parameters:

- `page`: Page number (1-based)
- `limit`: Number of items per page (default: 10, max: 100)

The response includes a pagination object:

```json
"pagination": {
  "total": 42,    // Total number of items
  "page": 1,      // Current page
  "limit": 10,    // Items per page
  "pages": 5      // Total number of pages
}
```

## Cross-References

- For authentication details, see the [Getting Started](getting_started.md) guide
- For architecture information, see the [Architecture](architecture.md) document
- For common questions, see the [FAQ](faq.md)
