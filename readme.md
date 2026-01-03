# 📚 Book Management API

A production-ready REST API for managing books and book borrowing system built with Go, implementing clean architecture principles.

## 🎯 Features

- **Authentication & Authorization**
  - JWT-based authentication
  - Protected endpoints with middleware
  - Password hashing with bcrypt

- **Book Management**
  - CRUD operations for books
  - Pagination support
  - Search by ISBN
  - Stock management

- **Borrow System**
  - Borrow books with automatic stock deduction
  - Return books with stock restoration
  - Transaction management with pessimistic locking
  - Borrow history tracking

- **Architecture**
  - Clean Architecture (Handler → Service → Repository)
  - Dependency Injection
  - Interface-based design for testability
  - Comprehensive unit tests
  - Graceful shutdown with signal handling
  - Interactive API documentation with Swagger

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Web Framework:** Chi Router
- **ORM:** GORM
- **Database:** PostgreSQL
- **Authentication:** JWT (golang-jwt/jwt)
- **Validation:** go-playground/validator
- **Testing:** testify
- **Configuration:** Viper
- **Documentation:** Swagger/OpenAPI

## 📁 Project Structure

```
book-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── database/                # Database connection & transaction manager
│   ├── models/                  # Data models
│   ├── repository/              # Data access layer
│   ├── services/                # Business logic layer
│   ├── handlers/                # HTTP handlers
│   ├── middlewares/             # HTTP middlewares
│   ├── routes/                  # Route definitions
│   └── utils/                   # Helper functions
├── docs/                        # Swagger documentation (auto-generated)
├── .env                         # Environment variables
├── go.mod                       # Go modules
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/Listantiyo/book-api.git
cd book-api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup database**
```bash
# Create database
psql -U postgres
CREATE DATABASE book_api;
\q
```

4. **Configure environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

`.env` example:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=book_api
DB_SSLMODE=disable

JWT_SECRET=your-super-secret-key-change-this
PORT=8080
```

5. **Generate Swagger documentation**
```bash
# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/server/main.go
```

6. **Run the application**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Accessing API Documentation

Once the server is running, you can access:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **API Base URL**: `http://localhost:8080/api/v1`
- **Health Check**: `http://localhost:8080/health`

The Swagger UI provides:
- Interactive API testing
- Complete endpoint documentation
- Request/response examples
- Authentication support (JWT Bearer token)
- Try-it-out functionality for all endpoints

## 📖 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Interactive Documentation

For complete, interactive API documentation with "Try it out" functionality, visit:
```
http://localhost:8080/swagger/index.html
```

### Quick Start Guide

#### 1. Register & Login via Swagger UI

1. Open `http://localhost:8080/swagger/index.html`
2. Find **POST /register** endpoint
3. Click "Try it out"
4. Enter registration details
5. Execute request
6. Use **POST /login** to get JWT token
7. Click "Authorize" button (top right)
8. Enter: `Bearer YOUR_TOKEN_HERE`
9. Now you can test protected endpoints!

---

### Manual API Testing

Below are curl examples for manual testing:

### Authentication Endpoints

#### Register User
```http
POST /register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### Book Endpoints

#### Get All Books (Public)
```http
GET /books?page=1&page_size=10
```

#### Get Book by ID (Public)
```http
GET /books/{id}
```

#### Create Book (Protected)
```http
POST /books
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "isbn": "9780132350884",
  "description": "A handbook of agile software craftsmanship",
  "stock": 10
}
```

#### Update Book (Protected)
```http
PUT /books/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Clean Code - Updated",
  "author": "Robert C. Martin",
  "isbn": "9780132350884",
  "description": "Updated description",
  "stock": 15
}
```

#### Delete Book (Protected)
```http
DELETE /books/{id}
Authorization: Bearer {token}
```

### Borrow Endpoints (All Protected)

#### Borrow Book
```http
POST /borrows
Authorization: Bearer {token}
Content-Type: application/json

{
  "book_id": 1
}
```

#### Return Book
```http
POST /borrows/return
Authorization: Bearer {token}
Content-Type: application/json

{
  "borrow_id": 1
}
```

#### Get My Borrows
```http
GET /borrows/me?page=1&page_size=10
Authorization: Bearer {token}
```

#### Get Borrow by ID
```http
GET /borrows/{id}
Authorization: Bearer {token}
```

## 🧪 Testing

Run all tests:
```bash
go test ./internal/services/... -v
```

Run tests with coverage:
```bash
go test ./internal/services/... -cover
```

Generate coverage report:
```bash
go test ./internal/services/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Regenerate Swagger Documentation

After modifying API endpoints or adding new handlers:

```bash
swag init -g cmd/server/main.go
```

This will update the `docs/` folder with latest API documentation.

## 🏗️ Architecture Decisions

### Clean Architecture Layers

1. **Handler Layer** - HTTP concerns (request/response, validation)
2. **Service Layer** - Business logic and orchestration
3. **Repository Layer** - Database operations and queries

### Key Patterns

- **Dependency Injection**: All dependencies injected via constructors
- **Interface Segregation**: Small, focused interfaces
- **Transaction Management**: Abstracted via TransactionManager interface
- **Pessimistic Locking**: `SELECT ... FOR UPDATE` for critical operations

### Why These Choices?

- **Testability**: Interfaces allow easy mocking
- **Maintainability**: Clear separation of concerns
- **Scalability**: Easy to add new features without breaking existing code
- **Safety**: Transaction + locking prevents race conditions

## 🔒 Security Features

- Password hashing with bcrypt (cost factor 10)
- JWT tokens with 24-hour expiration
- Protected endpoints via middleware
- SQL injection prevention (parameterized queries)
- Input validation on all endpoints

## 📈 Performance Considerations

- Pagination on list endpoints
- Database indexes on foreign keys
- Connection pooling (GORM default)
- Pessimistic locking only on critical paths
- Efficient query patterns (no N+1 queries)

## 🐛 Known Limitations

- No refresh token mechanism (JWT expires in 24h)
- No rate limiting implemented
- No caching layer
- Pessimistic locking may cause performance bottleneck under high concurrency

## 🔮 Future Improvements

- [ ] Add refresh token support
- [ ] Implement Redis caching for book list
- [ ] Add rate limiting middleware
- [ ] Implement role-based access control (Admin/User)
- [ ] Add integration tests
- [ ] Add Docker support
- [ ] CI/CD pipeline setup
- [ ] Request timeout with context propagation
- [ ] Distributed tracing
- [ ] WebSocket support for real-time notifications
- [ ] File upload for book covers

## 📝 License

MIT License - feel free to use this project for learning or commercial purposes.

## 👤 Author

Listantiyo - [GitHub](https://github.com/Listantiyo)

## 🙏 Acknowledgments

- Inspired by clean architecture principles
- Built as a learning project for Go backend development
- Swagger documentation generated with swaggo/swag
- Special thanks to the Go community

---

## 📚 Additional Resources

- **Swagger UI**: `http://localhost:8080/swagger/index.html` (when server is running)
- **Go Documentation**: https://go.dev/doc/
- **GORM Documentation**: https://gorm.io/docs/
- **Chi Router**: https://github.com/go-chi/chi
- **Swaggo**: https://github.com/swaggo/swag
