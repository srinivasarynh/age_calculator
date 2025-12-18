# User API - Go Backend with Age Calculation

A RESTful API built with Go that manages users with dynamic age calculation based on date of birth.

## Features

✅ **Complete CRUD operations** for users  
✅ **Dynamic age calculation** from DOB  
✅ **Pagination support** for listing users  
✅ **Input validation** with go-playground/validator  
✅ **Structured logging** with Uber Zap  
✅ **Request ID tracking** for debugging  
✅ **Request duration logging**  
✅ **Docker support** for easy deployment  
✅ **Clean architecture** with layered structure  
✅ **Unit tests** for age calculation  

## Tech Stack

- **Framework**: GoFiber
- **Database**: PostgreSQL with raw SQL queries
- **Validation**: go-playground/validator
- **Logging**: Uber Zap
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── config/
│   └── config.go                   # Configuration management
├── db/
│   ├── migrations/
│   │   └── 001_create_users.sql   # Database schema
│   └── queries/
│       └── users.sql               # SQL queries for SQLC
├── internal/
│   ├── handler/
│   │   └── user_handler.go        # HTTP handlers
│   ├── repository/
│   │   └── user_repository.go     # Database operations
│   ├── service/
│   │   ├── user_service.go        # Business logic
│   │   └── user_service_test.go   # Unit tests
│   ├── routes/
│   │   └── routes.go               # Route definitions
│   ├── middleware/
│   │   └── middleware.go           # Custom middleware
│   ├── models/
│   │   └── user.go                 # Data models
│   └── logger/
│       └── logger.go               # Logger configuration
├── docker-compose.yml              # Docker services
├── Dockerfile                      # Application container
├── Makefile                        # Build automation
├── go.mod                          # Go dependencies
└── README.md                       # This file
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+ (or use Docker)
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd user-api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Start PostgreSQL** (if not using Docker)
```bash
# Using Docker:
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=userdb \
  -p 5432:5432 \
  postgres:15-alpine
```

5. **Run database migrations**
```bash
# Apply the schema manually:
psql -U postgres -d userdb -f db/migrations/001_create_users.sql

# Or use golang-migrate:
make migrate-up
```

6. **Run the application**
```bash
make run
# Or directly:
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

### Using Docker

**Start everything with Docker Compose:**
```bash
make docker-up
```

**Stop services:**
```bash
make docker-down
```

**View logs:**
```bash
make docker-logs
```

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check
```http
GET /health
```

### 1. Create User
```http
POST /api/v1/users
Content-Type: application/json

{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### 2. Get User by ID
```http
GET /api/v1/users/1
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 34
}
```

### 3. List All Users (with Pagination)
```http
GET /api/v1/users?page=1&page_size=10
```

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": 1,
      "name": "Alice",
      "dob": "1990-05-10",
      "age": 34
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

### 4. Update User
```http
PUT /api/v1/users/1
Content-Type: application/json

{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### 5. Delete User
```http
DELETE /api/v1/users/1
```

**Response: 204 No Content**

## Testing

### Run all tests
```bash
make test
```

### Run tests with coverage
```bash
make test-coverage
```

### Test the age calculation function
```bash
go test -v ./internal/service -run TestCalculateAge
```

## Validation Rules

### Create/Update User Request
- **name**: Required, minimum 2 characters, maximum 100 characters
- **dob**: Required, must be in format `YYYY-MM-DD`

## Error Responses

All error responses follow this format:
```json
{
  "error": "Error message",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `204` - No Content (successful deletion)
- `400` - Bad Request (validation error)
- `404` - Not Found
- `500` - Internal Server Error

## Middleware Features

### 1. Request ID
Every request gets a unique ID that's:
- Returned in `X-Request-ID` header
- Included in all log entries
- Included in error responses

### 2. Request Logger
Logs all requests with:
- Request ID
- HTTP method and path
- Response status code
- Request duration
- Client IP
- User agent

Example log:
```json
{
  "level": "info",
  "ts": "2025-12-17T10:30:45.123Z",
  "msg": "HTTP Request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "GET",
  "path": "/api/v1/users/1",
  "status": 200,
  "duration": "2.5ms",
  "ip": "127.0.0.1"
}
```

## Age Calculation Logic

The age is calculated dynamically using Go's `time` package:

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

## Make Commands

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage report
make clean          # Clean build artifacts
make deps           # Download dependencies
make docker-build   # Build Docker image
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make docker-logs    # View Docker logs
```

## Example Usage with curl

### Create a user
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "dob": "1995-08-15"
  }'
```

### Get a user
```bash
curl http://localhost:8080/api/v1/users/1
```

### List users with pagination
```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=5"
```

### Update a user
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "dob": "1995-08-20"
  }'
```

### Delete a user
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    dob DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## License

MIT License
