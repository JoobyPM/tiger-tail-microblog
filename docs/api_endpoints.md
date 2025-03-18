# API Endpoints

This document describes the API endpoints available in the Tiger-Tail Microblog application.

## Authentication

### POST /api/auth/login

Authenticates a user and returns a JWT token.

**Request:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "token": "string",
  "expires_at": "string (ISO 8601 datetime)"
}
```

### POST /api/auth/register

Registers a new user.

**Request:**
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "id": "string",
  "username": "string",
  "created_at": "string (ISO 8601 datetime)"
}
```

## Posts

### GET /api/posts

Returns a list of posts.

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of posts per page (default: 10)
- `user_id`: Filter by user ID (optional)

**Response:**
```json
{
  "posts": [
    {
      "id": "string",
      "content": "string",
      "user_id": "string",
      "username": "string",
      "created_at": "string (ISO 8601 datetime)",
      "updated_at": "string (ISO 8601 datetime)"
    }
  ],
  "pagination": {
    "total": "integer",
    "page": "integer",
    "limit": "integer",
    "pages": "integer"
  }
}
```

### POST /api/posts

Creates a new post.

**Request:**
```json
{
  "content": "string"
}
```

**Response:**
```json
{
  "id": "string",
  "content": "string",
  "user_id": "string",
  "username": "string",
  "created_at": "string (ISO 8601 datetime)",
  "updated_at": "string (ISO 8601 datetime)"
}
```

### GET /api/posts/{id}

Returns a specific post.

**Response:**
```json
{
  "id": "string",
  "content": "string",
  "user_id": "string",
  "username": "string",
  "created_at": "string (ISO 8601 datetime)",
  "updated_at": "string (ISO 8601 datetime)"
}
```

### PUT /api/posts/{id}

Updates a specific post.

**Request:**
```json
{
  "content": "string"
}
```

**Response:**
```json
{
  "id": "string",
  "content": "string",
  "user_id": "string",
  "username": "string",
  "created_at": "string (ISO 8601 datetime)",
  "updated_at": "string (ISO 8601 datetime)"
}
```

### DELETE /api/posts/{id}

Deletes a specific post.

**Response:**
```json
{
  "success": true
}
```

## Users

### GET /api/users/{id}

Returns a specific user.

**Response:**
```json
{
  "id": "string",
  "username": "string",
  "bio": "string",
  "created_at": "string (ISO 8601 datetime)"
}
```

### PUT /api/users/{id}

Updates a specific user.

**Request:**
```json
{
  "bio": "string"
}
```

**Response:**
```json
{
  "id": "string",
  "username": "string",
  "bio": "string",
  "created_at": "string (ISO 8601 datetime)",
  "updated_at": "string (ISO 8601 datetime)"
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request

```json
{
  "error": "string",
  "message": "string"
}
```

### 401 Unauthorized

```json
{
  "error": "string",
  "message": "string"
}
```

### 403 Forbidden

```json
{
  "error": "string",
  "message": "string"
}
```

### 404 Not Found

```json
{
  "error": "string",
  "message": "string"
}
```

### 500 Internal Server Error

```json
{
  "error": "string",
  "message": "string"
}
```
