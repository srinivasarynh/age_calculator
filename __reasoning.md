# Architecture & Design Decisions

## Overview
This document explains the key architectural decisions, design patterns, and trade-offs made during the development of the User API.

---

## 1. Architecture Pattern: Clean Architecture

### Decision
Implemented a layered architecture with clear separation of concerns:
```
Handler → Service → Repository → Database
```

### Reasoning
- **Maintainability**: Each layer has a single responsibility, making code easier to understand and modify
- **Testability**: Layers can be tested independently with mock implementations
- **Scalability**: Easy to swap implementations (e.g., change from PostgreSQL to MongoDB)
- **Dependency Rule**: Inner layers (service, repository) don't depend on outer layers (handler)

---

## 2. Database Layer: Raw SQL vs ORM

### Decision
Used **raw SQL queries** with `database/sql` package instead of an ORM (GORM, Ent).

### Reasoning
- **Performance**: Direct SQL is faster with no ORM overhead
- **Control**: Full control over queries, joins, and optimizations
- **Transparency**: Easy to see exactly what queries are executed
- **SQLC Ready**: Prepared for SQLC code generation (included in structure)
- **Learning**: Better understanding of SQL for developers

### Trade-offs
- **More Code**: Need to write SQL queries manually and handle scanning
- **No Migrations**: No built-in migration tool (but SQL migrations are standard)
- **Type Safety**: Less compile-time safety compared to type-safe ORMs like Ent

### Alternative Considered
Could use SQLC for type-safe generated code, but kept it simple with raw SQL to demonstrate fundamentals.

---

## 3. Web Framework: GoFiber

### Decision
Chose **GoFiber** over standard library or alternatives (Gin, Echo, Chi).

### Reasoning
- **Performance**: Built on Fasthttp, one of the fastest Go web frameworks
- **Express-like API**: Familiar syntax for developers coming from Node.js
- **Rich Middleware**: Built-in middleware for CORS, recovery, compression
- **Modern**: Active development and strong community support
- **Low Memory**: Efficient memory usage compared to net/http
---

## 4. Validation: go-playground/validator

### Decision
Used **go-playground/validator** for input validation with struct tags.

### Reasoning
- **Declarative**: Validation rules defined directly on structs
- **Standard**: Most popular validation library in Go ecosystem
- **Rich Rules**: Extensive built-in validators (min, max, email, datetime, etc.)
- **Custom Validators**: Easy to add custom validation rules
- **Clear Errors**: Provides detailed validation error messages

### Example
```go
type CreateUserRequest struct {
    Name string `json:"name" validate:"required,min=2,max=100"`
    DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}
```

### Alternative Considered
Could validate manually in handlers, but struct tags are cleaner and more maintainable.

---

## 5. Logging: Uber Zap

### Decision
Implemented **Uber Zap** for structured logging.

### Reasoning
- **Performance**: Fastest structured logger for Go (allocation-free)
- **Structured**: JSON output perfect for log aggregation (ELK, Datadog)
- **Type-Safe**: Compile-time type checking for log fields
- **Production-Ready**: Battle-tested at Uber and other companies
- **Development Mode**: Beautiful console output during development

### Example Log Output
```json
{
  "level": "info",
  "ts": "2025-12-17T10:30:45.123Z",
  "msg": "HTTP Request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "GET",
  "path": "/api/v1/users/1",
  "status": 200,
  "duration": "2.5ms"
}
```

### Why Not Standard Logger?
Standard `log` package lacks structured logging, making it hard to parse and analyze logs at scale.

---

## 6. Age Calculation: Dynamic Computation

### Decision
Calculate age **dynamically** in the service layer, not stored in database.

### Reasoning
- **Always Accurate**: Age is always current, no stale data
- **Single Source of Truth**: DOB is the only stored value
- **Simple Logic**: Calculation is lightweight and fast
- **No Updates Needed**: Don't need to update ages periodically

### Implementation
```go
func CalculateAge(dob time.Time) int {
    now := time.Now()
    age := now.Year() - dob.Year()
    
    // Adjust if birthday hasn't occurred this year
    if now.Month() < dob.Month() || 
       (now.Month() == dob.Month() && now.Day() < dob.Day()) {
        age--
    }
    
    return age
}
```

### Edge Cases Handled
- Leap year birthdays
- Birthday today
- Birthday tomorrow/yesterday
- Month/day boundary conditions

### Alternative Considered
Could store age in database and update via cron job, but this adds complexity and potential inconsistency.

---

## 7. Error Handling Strategy

### Decision
Implemented **custom error types** and **global error handler**.

### Reasoning
- **Consistency**: All errors follow the same format
- **Request Tracking**: Every error includes request ID for debugging
- **Clean Handlers**: Handlers don't need to handle every error case
- **HTTP Mapping**: Automatic mapping of errors to HTTP status codes

### Error Flow
```
Service Error → Handler Check → Global Handler → JSON Response
```

### Example Error Response
```json
{
  "error": "User not found",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 8. Middleware Design

### Decision
Implemented three custom middleware: **RequestID**, **Logger**, and **ErrorHandler**.

### RequestID Middleware
**Purpose**: Add unique ID to every request
- Enables request tracing across logs
- Helps debug issues in production
- Can be used for distributed tracing later

### Logger Middleware
**Purpose**: Log all HTTP requests with metrics
- Request method, path, status code
- Duration for performance monitoring
- IP and user agent for security
- All logs include request ID

### ErrorHandler Middleware
**Purpose**: Centralized error handling
- Converts errors to JSON responses
- Maps error types to HTTP status codes
- Includes request ID in error responses

---

## 9. Pagination Implementation

### Decision
Added **offset-based pagination** with configurable page size.

### Reasoning
- **Simple**: Easy to implement and understand
- **Standard**: Most common pagination approach
- **Flexible**: Users can control page size (1-100)
- **Metadata**: Returns total count, pages, current page

### Response Structure
```json
{
  "users": [...],
  "total": 150,
  "page": 2,
  "page_size": 10,
  "total_pages": 15
}
```

### Query Parameters
- `page`: Current page number (default: 1)
- `page_size`: Items per page (default: 10, max: 100)

### Alternative Considered
Could use **cursor-based pagination** for better performance with large datasets, but offset-based is simpler and sufficient for most use cases.

---

## 10. Docker Setup

### Decision
Provided **multi-stage Dockerfile** and **Docker Compose** for full stack.

### Multi-Stage Build
```dockerfile
Stage 1: Builder (golang:1.21-alpine)
  → Install dependencies
  → Build binary

Stage 2: Runtime (alpine:latest)
  → Copy binary only
  → Minimal image size
```

### Benefits
- **Small Images**: Final image ~15MB vs ~800MB with full Go image
- **Security**: Fewer dependencies = smaller attack surface
- **Fast Deploys**: Smaller images deploy faster

### Docker Compose
Orchestrates two services:
1. **PostgreSQL**: Database with health check
2. **API**: Application that waits for DB to be ready

### Why Docker?
- **Consistency**: Same environment across development and production
- **Easy Setup**: `make docker-up` starts everything
- **Isolation**: No conflicts with local installations

---

## 11. Configuration Management

### Decision
Used **environment variables** with sensible defaults.

### Reasoning
- **12-Factor App**: Follows cloud-native principles
- **Flexibility**: Easy to change config without rebuilding
- **Security**: Secrets not hardcoded in source
- **Docker-Friendly**: Easy to override in containers

### Configuration Sources
1. Environment variables (highest priority)
2. .env file (development)
3. Default values (fallback)

---

## 12. Testing Strategy

### Decision
Implemented **unit tests** for age calculation, with structure for more tests.

### Current Tests
- Age calculation with various DOBs
- Edge cases (birthday today, leap years)
- Boundary conditions

### Future Testing Recommendations
1. **Unit Tests**: Service layer business logic
2. **Integration Tests**: Repository with test database
3. **E2E Tests**: Full API flows with test containers
4. **Benchmark Tests**: Performance testing for age calculation

### Why Start with Unit Tests?
- Fastest to run
- No external dependencies
- Demonstrate testing approach
- Most critical business logic

---

## 13. Date Format: ISO 8601

### Decision
Used **YYYY-MM-DD** format for dates.

### Reasoning
- **ISO 8601**: International standard
- **Sortable**: Lexicographic sorting works
- **Unambiguous**: No confusion between MM/DD vs DD/MM
- **Database Native**: PostgreSQL DATE type expects this format
- **JSON Standard**: Widely used in REST APIs

---

## 14. API Versioning

### Decision
Included **v1** in API path: `/api/v1/users`.

### Reasoning
- **Future-Proof**: Can introduce v2 without breaking v1
- **Clear Intent**: Users know which version they're using
- **Best Practice**: Standard in REST API design
- **Migration Path**: Easier to deprecate old versions

---

## 15. HTTP Status Codes

### Decision
Used **semantic HTTP status codes** consistently.

### Mapping
| Status | Use Case |
|--------|----------|
| 200 OK | Successful GET/PUT |
| 201 Created | Successful POST |
| 204 No Content | Successful DELETE |
| 400 Bad Request | Validation errors |
| 404 Not Found | Resource doesn't exist |
| 500 Internal Server Error | Server errors |

### Why Not Just 200?
Proper status codes:
- Help clients handle errors correctly
- Follow REST conventions
- Enable better API monitoring
- Improve developer experience

---

## 16. Database Connection Management

### Decision
Used **connection pooling** with configured limits.

### Configuration
```go
db.SetMaxOpenConns(25)  // Max concurrent connections
db.SetMaxIdleConns(5)   // Keep 5 connections ready
```

### Reasoning
- **Performance**: Reuse connections instead of creating new ones
- **Resource Management**: Prevent connection exhaustion
- **Concurrency**: Handle multiple requests efficiently
- **Production-Ready**: Essential for high-traffic scenarios

---

## 17. Code Organization Principles

### Package Structure
- **`cmd/`**: Application entrypoints (can have multiple)
- **`config/`**: Configuration loading and management
- **`internal/`**: Private application code (not importable)
- **`db/`**: Database migrations and queries

### Internal Package Benefits
- Prevents external packages from importing internal code
- Enforces API boundaries
- Common Go convention

### Why This Structure?
- **Standard Layout**: Follows golang-standards/project-layout
- **Scalable**: Easy to add new services/handlers
- **Clear Boundaries**: Each package has clear responsibility
- **Team-Friendly**: New developers can find code easily

---

## 18. Graceful Shutdown

### Decision
Implemented **graceful shutdown** with signal handling.

### Mechanism
```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

<-quit // Wait for signal
app.ShutdownWithContext(context.Background())
```

### Why Important?
- **Complete Requests**: Finish processing current requests
- **Close Connections**: Properly close DB connections
- **Data Integrity**: Prevent data corruption
- **Cloud-Ready**: Required for Kubernetes rolling updates

---

## 19. Performance Considerations

### Optimizations Made
1. **Indexed DOB Field**: Faster queries on date ranges
2. **Connection Pooling**: Reuse database connections
3. **No N+1 Queries**: Single query for list operations
4. **Efficient Scanning**: Direct struct scanning
5. **Minimal Allocations**: Zap logger is allocation-free

### Potential Improvements
- Add Redis caching for frequently accessed users
- Implement batch operations for bulk creates
- Use prepared statements for repeated queries
- Add database read replicas for scaling reads

---

## 20. Security Considerations

### Current Implementation
✅ SQL injection prevention (parameterized queries)  
✅ Input validation (length, format)  
✅ Error message sanitization (no internal details exposed)  
✅ CORS middleware included  
✅ Recovery middleware (prevents crashes)  

### Production Recommendations
- [ ] Add authentication (JWT, OAuth2)
- [ ] Implement rate limiting per IP
- [ ] Add HTTPS/TLS support
- [ ] Implement authorization/RBAC
- [ ] Add request size limits
- [ ] Enable security headers (helmet)
- [ ] Add API key management
- [ ] Implement audit logging

---

## 21. Monitoring & Observability

### Current Implementation
- Structured logging with Zap
- Request ID tracing
- Request duration metrics
- HTTP status code tracking

### Production Recommendations
- **Metrics**: Prometheus metrics endpoint
- **Tracing**: OpenTelemetry integration
- **Alerts**: Set up alerts on error rates
- **Dashboards**: Grafana for visualization
- **APM**: Application Performance Monitoring (Datadog, New Relic)

---

## 22. Trade-offs Summary

### What We Optimized For
1. **Code Quality**: Clean, maintainable, well-documented
2. **Developer Experience**: Easy to understand and modify
3. **Production-Ready**: Logging, errors, Docker, graceful shutdown
4. **Performance**: Fast framework, connection pooling, efficient queries

### What We Didn't Optimize For
1. **Minimal Code**: Chose clarity over brevity
2. **Cutting Edge**: Used stable, proven technologies
3. **Every Feature**: Focused on core requirements + best practices

---

## 23. Future Enhancements

### High Priority
- [ ] Authentication & authorization
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Comprehensive test suite
- [ ] CI/CD pipeline

### Medium Priority
- [ ] Rate limiting
- [ ] Caching layer (Redis)
- [ ] Background jobs (for reports, emails)
- [ ] Soft deletes

### Nice to Have
- [ ] GraphQL endpoint
- [ ] WebSocket support
- [ ] Multiple database support
- [ ] Multi-tenancy

---

## 24. Lessons & Best Practices

### Key Takeaways
1. **Start Simple**: Begin with working code, refactor later
2. **Layer Properly**: Separation of concerns pays off
3. **Log Everything**: Structured logs are invaluable
4. **Test Core Logic**: Unit tests for business logic first
5. **Document Decisions**: This file helps future developers

### Patterns to Continue
- Dependency injection (pass dependencies explicitly)
- Interface-driven design (easy to mock/swap)
- Error wrapping (preserve context up the stack)
- Consistent naming (similar functions across layers)

### Patterns to Avoid
- Global state (use dependency injection)
- Magic values (use constants)
- God objects (keep packages focused)
- Implicit behavior (make it explicit)

---

## Conclusion

This API demonstrates production-ready Go development with:
- Clean architecture and separation of concerns
- Proper error handling and logging
- Comprehensive documentation
- Docker deployment support
- Extensibility for future enhancements

The design prioritizes **maintainability**, **testability**, and **developer experience** while maintaining high performance and production readiness.

---

**Author**: Built following Go best practices and industry standards  
**Date**: December 2025  
**Version**: 1.0
