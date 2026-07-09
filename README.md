# RiffLog

RiffLog is a RESTful backend API for tracking guitar practice sessions. It allows users to create an account, authenticate using JSON Web Tokens (JWT), record practice sessions, browse available practice skills, and view statistics about their practice history.

This project was developed as a portfolio piece to demonstrate modern backend development practices using Go. Rather than focusing on building a production-ready application, the goal was to showcase clean architecture, REST API design, authentication, PostgreSQL integration, Docker-based development, and automated testing.

---

## Features

* User registration and authentication
* JWT-based authentication and authorization
* Password hashing with bcrypt
* Browse available practice skills
* Create, retrieve, update, and delete practice sessions
* Filter practice sessions by skill and date range
* Aggregate practice statistics, including:

  * Total practice time
  * Total practice sessions
  * Most practiced skill
  * Longest practice session
* PostgreSQL persistence
* SQL database migrations
* Docker development environment
* Unit and integration test coverage

---

## Technology Stack

| Technology              | Purpose                       |
| ----------------------- | ----------------------------- |
| Go                      | Backend language              |
| Gin                     | HTTP routing and middleware   |
| PostgreSQL              | Relational database           |
| pgx                     | PostgreSQL driver             |
| Docker & Docker Compose | Local development environment |
| JWT                     | Authentication                |
| bcrypt                  | Password hashing              |
| golang-migrate          | Database migrations           |

---

## Project Architecture

The application follows a layered architecture to separate HTTP concerns, business logic, and data persistence.

```text
HTTP Request
      │
      ▼
Gin Router
      │
      ▼
Middleware (Authentication)
      │
      ▼
Handler
      │
      ▼
Service
      │
      ▼
Repository
      │
      ▼
PostgreSQL
```

Responsibilities are separated into:

* **Handlers** – Parse HTTP requests and build HTTP responses.
* **Services** – Implement business rules and validation.
* **Repositories** – Execute SQL queries and map database results.
* **Middleware** – Authenticate requests and populate the authenticated user context.

---

## Project Structure

```text
cmd/
internal/
    auth/
    bootstrap/
    config/
    database/
    handlers/
    middleware/
    models/
    repository/
    services/
migrations/
docs/
```

---

## Getting Started

### Prerequisites

* Go
* Docker Desktop
* PostgreSQL (via Docker Compose)
* golang-migrate CLI

### About development workflow

During development, PostgreSQL runs in Docker while the Go API is run directly from the local development environment. This provides a consistent database environment while allowing fast compilation, debugging, and testing of the Go application.

### Clone the repository

```bash
git clone https://github.com/thetramp22/rifflog.git
cd rifflog
```

### Configure environment variables

```bash
cp .env.example .env
```

Update the values in `.env` to match your local environment.

### Start the database

```bash
docker compose up --build -d
```

### Run database migrations

```bash
migrate -path migrations \
-database "postgres://rifflog:devriffs@localhost:5433/rifflog?sslmode=disable" up
```

### Start the API

```bash
go run ./cmd/api
```

---

## Running Tests

Run all tests:

```bash
go test ./...
```

---

## API Documentation

Complete API documentation is available in:

```text
docs/api.md
```

The documentation includes:

* Request and response examples
* Authentication requirements
* Query parameters
* Path parameters
* Error responses

---

## Design Decisions

Several design decisions were intentionally made while developing this project:

* JWT authentication is handled through dedicated middleware rather than requiring client-supplied user IDs.
* Repository methods are responsible for translating database-specific behavior into application-level errors.
* Context is propagated through the service and repository layers using Go's standard `context.Context`.
* Practice session ownership is enforced at the database query level to prevent users from accessing or modifying another user's data.
* SQL queries are written directly rather than using an ORM to demonstrate familiarity with relational database design and PostgreSQL.

---

## Future Improvements

Possible future enhancements include:

* Refresh token support
* Password reset workflow
* OpenAPI (Swagger) documentation
* CI/CD pipeline
* Pagination for large result sets
* Practice goals and streak tracking
* Frontend client application

---

## What I Learned

Building RiffLog provided practical experience with:

* Designing RESTful APIs in Go
* Layered application architecture
* PostgreSQL schema design and SQL
* JWT authentication and authorization
* Writing middleware with Gin
* Integration and unit testing
* Docker-based local development
* Database migrations
* Structuring a maintainable Go project

---

## License

This project is available under the MIT License.
