# Tiger-Tail Microblog Architecture

## Overview

This document outlines the architecture of the Tiger-Tail Microblog application, following our [Tiger Style](tiger_style.md) principles of safety, performance, and developer experience.

## System Components

Tiger-Tail follows a clean, layered architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Server Layer                      │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │   Handlers      │  │   Middleware    │  │   Router    │  │
│  └─────────────────┘  └─────────────────┘  └─────────────┘  │
└───────────────────────────────┬─────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                          │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │  Post Service   │  │  User Service   │                   │
│  └─────────────────┘  └─────────────────┘                   │
└───────────────────────────────┬─────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                           │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │   Post Entity   │  │   User Entity   │                   │
│  └─────────────────┘  └─────────────────┘                   │
└───────────────────────────────┬─────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      Data Access Layer                      │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │ Post Repository │  │ User Repository │                   │
│  └─────────────────┘  └─────────────────┘                   │
└───────────────────────────────┬─────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      Storage Layer                          │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │    PostgreSQL   │  │      Redis      │                   │
│  └─────────────────┘  └─────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

### HTTP Server Layer

The HTTP server layer handles incoming HTTP requests, routing, and middleware. It's implemented using the Gin framework.

### Service Layer

The service layer contains the business logic of the application. It orchestrates the flow of data between the HTTP layer and the data access layer.

### Domain Layer

The domain layer defines the core entities and business rules. It's the heart of the application and is independent of any external frameworks or libraries.

### Data Access Layer

The data access layer handles persistence and retrieval of data. It abstracts away the details of the underlying storage mechanisms.

### Storage Layer

The storage layer consists of PostgreSQL for persistent storage and Redis for caching.

## Data Flow

The data flow in Tiger-Tail follows a clean, predictable pattern:

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │────▶│  Server  │────▶│  Service │────▶│Repository│
│          │◀────│          │◀────│          │◀────│          │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                                                        │
                                                        ▼
                                                   ┌──────────┐
                                                   │  Redis   │
                                                   │  Cache   │
                                                   └──────────┘
                                                        │
                                                        ▼
                                                   ┌──────────┐
                                                   │PostgreSQL│
                                                   │ Database │
                                                   └──────────┘
```

1. Client sends a request to the server
2. Server routes the request to the appropriate handler
3. Handler calls the service layer
4. Service layer applies business logic
5. Service layer calls the repository layer
6. Repository layer checks Redis cache first
7. If data is not in cache, repository fetches from PostgreSQL
8. Repository updates Redis cache with fetched data
9. Data flows back up through the layers to the client

## Performance Considerations

Following our Tiger Style principles, we've optimized for performance:

1. **Redis Caching**: Frequently accessed data is cached in Redis to reduce database load
2. **Connection Pooling**: Database connections are pooled for efficient reuse
3. **Pagination**: API endpoints that return lists support pagination to limit response size
4. **Indexing**: Database tables are properly indexed for fast queries
5. **Concurrency**: Go's goroutines are used for concurrent processing where appropriate

## Security Considerations

Security is our top priority:

1. **Input Validation**: All user input is validated before processing
2. **Parameterized Queries**: SQL queries use parameterization to prevent injection attacks
3. **Authentication**: API endpoints are protected with authentication where appropriate
4. **Rate Limiting**: API endpoints are rate-limited to prevent abuse
5. **HTTPS**: All production deployments require HTTPS
6. **Error Handling**: Errors are logged but not exposed to clients in production
