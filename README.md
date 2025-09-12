# Graphy - GraphQL API Server

A production-ready GraphQL API server built with Go using gqlgen, featuring JWT authentication, PostgreSQL database, and comprehensive note management. This project demonstrates clean architecture principles with domain-driven design.

## Features

- 🚀 **GraphQL API** - Built with gqlgen for type-safe GraphQL operations
- 🔐 **JWT Authentication** - Secure user authentication with Bearer token support
- 📝 **Note Management** - Full CRUD operations for user notes with authorization
- 👤 **User Management** - User registration, login, and profile management
- 🐘 **PostgreSQL Database** - Robust data persistence with connection pooling (pgx/v5)
- 🔄 **Database Migrations** - Automated schema migrations using golang-migrate
- 🏗️ **Clean Architecture** - Domain-driven design with repository pattern
- 🔧 **Dependency Injection** - Centralized dependency management
- 🐳 **Docker Support** - Containerized PostgreSQL setup
- � **Hot Reload** - Development environment with Air
- 🎨 **Colored Logging** - Enhanced logging with slog and colored output
- 🌐 **CORS Support** - Configurable cross-origin resource sharing
- ⚡ **Query Caching** - Built-in GraphQL query caching and automatic persisted queries

## Architecture Overview

This project follows **Clean Architecture** principles:

- **Domain Layer** (`internal/domain/`) - Business logic and entities
- **Infrastructure Layer** (`internal/infrastructure/`) - External concerns (database)
- **Application Layer** (`graph/`) - GraphQL resolvers and schema
- **Configuration Layer** (`config/`) - Application configuration and setup

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.24.6 or later)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Air](https://github.com/air-verse/air) (for hot reload during development)

```bash
# Install Air globally for hot reload
go install github.com/air-verse/air@latest
```

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/daffadon/graphy.git
cd graphy
```

### 2. Environment Setup

Create your local configuration file:

```bash
cp config.example.yaml config.local.yaml
```

Edit `config.local.yaml` with your database configuration:

```yaml
database:
  sql:
    protocol: postgres
    host: localhost
    user: postgres
    password: your_password
    port: 5432
    name: graphy_db
    sslmode: disable

app:
  port: 8080

server:
  cors:
    allow_origins: "*"
    allow_methods: "GET, POST, PUT, DELETE, OPTIONS, PATCH"
    allow_headers: "Content-Type, Authorization, X-Requested-With, X-CSRF-Token, Accept, Origin, Cache-Control, X-File-Name, X-File-Type, X-File-Size"
    expose_headers: "Content-Length, Content-Range"
    max_age: 86400
    allow_credential: true
```

### 3. Database Setup

Start PostgreSQL using Docker Compose:

```bash
# Create .env file for Docker Compose (if it doesn't exist)
echo "DB_USER=postgres" > .env
echo "DB_PASSWORD=your_password" >> .env
echo "DB_NAME=graphy_db" >> .env

# Start PostgreSQL container
docker-compose up -d postgresql
```

Verify the database is running:

```bash
docker ps
```

You should see the PostgreSQL container running on port 5432.

### 4. GraphQL Setup and Code Generation

#### First Time Setup (Initialize GraphQL)

If this is your first time setting up the project or if `tools.go` doesn't exist:

```bash
# Initialize gqlgen (creates tools.go and sets up GraphQL)
make gql-init
```

This command will:

- Create `tools.go` with necessary gqlgen imports
- Run `go mod tidy` to download dependencies
- Initialize gqlgen with basic GraphQL setup

#### Generate GraphQL Code

After modifying the GraphQL schema (`graph/schema.graphqls`), regenerate the code:

```bash
# Generate GraphQL resolvers and models
make gql-generate
```

This command generates:

- `graph/generated.go` - GraphQL server implementation
- `graph/model/models_gen.go` - GraphQL models based on schema
- `graph/schema.resolvers.go` - Resolver implementations (preserves existing code)

### 5. Install Dependencies

```bash
go mod download
go mod tidy
```

### 6. Database Migration

The application will automatically run migrations on startup. The migration system will create:

- `users` table (id, email, fullname, password, timestamps)
- `notes` table (id, title, description, text, user_id, timestamps with foreign key)

Migration files are located in `internal/pkg/database/migrations/postgresql/`.

### 7. Run the Application

#### Development Mode (Recommended)

Use Air for hot reload during development:

```bash
air
```

Air configuration is in `.air.toml` and will:

- Watch for file changes in `.go`, `.yaml`, `.html` files
- Automatically rebuild and restart the server
- Log build errors to `tmp/build-errors.log`

#### Production Mode

```bash
go run server.go
```

The server will start on `http://localhost:8080` (or the port specified in your config).

### 8. Access GraphQL Playground

Once the server is running, access the GraphQL playground at:

```url
http://localhost:8080/gq
```

The playground provides:

- Schema exploration
- Query/mutation testing
- Real-time documentation

## Project Structure

```text
graphy/
├── cmd/                          # Application commands and bootstrapping
│   └── bootstrap.go             # Dependency injection and app initialization
├── config/                      # Configuration packages
│   ├── database/               # Database connection configuration
│   │   └── postgresql.go       # PostgreSQL connection setup
│   ├── env/                    # Environment configuration loader
│   │   └── env.go              # Viper configuration management
│   ├── logger/                 # Logging configuration
│   │   └── slog.go             # Structured logging with colors
│   └── router/                 # HTTP router configuration
│       └── http.go             # Chi router with CORS and middleware
├── graph/                       # GraphQL layer
│   ├── generated.go            # Generated GraphQL server code (auto-generated)
│   ├── resolver.go             # Main resolver with dependencies
│   ├── schema.graphqls         # GraphQL schema definition
│   ├── schema.resolvers.go     # Resolver implementations
│   └── model/                  # GraphQL models
│       └── models_gen.go       # Generated models from schema
├── internal/                    # Private application code
│   ├── domain/                 # Business logic layer
│   │   ├── auth/              # Authentication middleware
│   │   │   └── middleware.go  # JWT authentication middleware
│   │   ├── dto/               # Data transfer objects
│   │   │   └── model.go       # Internal data models
│   │   ├── notes/             # Notes domain logic
│   │   │   └── notes.go       # Note repository and business logic
│   │   └── users/             # Users domain logic
│   │       └── users.go       # User repository and business logic
│   ├── infrastructure/         # External concerns
│   │   └── database/          # Database infrastructure
│   │       └── querier.go     # Database query interface
│   └── pkg/                   # Internal packages
│       ├── database/          # Database utilities
│       │   ├── migrations/    # Database migration files
│       │   │   └── postgresql/
│       │   │       ├── 000001_create_users_table.up.sql
│       │   │       ├── 000001_create_users_table.down.sql
│       │   │       ├── 000002_create_notes_table.up.sql
│       │   │       └── 000002_create_notes_table.down.sql
│       │   └── postgresql/    # PostgreSQL utilities
│       │       └── postgresql.go # Migration runner
│       ├── jwt/               # JWT utilities
│       │   └── jwt.go         # Token generation and validation
│       └── utils/             # Common utilities
│           └── password.go    # Password hashing utilities
├── tmp/                        # Temporary build files
│   ├── build-errors.log       # Air build error logs
│   └── main.exe               # Compiled binary
├── .air.toml                   # Air hot reload configuration
├── config.example.yaml         # Configuration template
├── config.local.yaml           # Local configuration (create this)
├── docker-compose.yml          # Docker services configuration
├── gqlgen.yml                  # GraphQL generation configuration
├── go.mod                      # Go module dependencies
├── go.sum                      # Go module checksums
├── Makefile                    # Build and setup commands
├── server.go                   # Application entry point
└── tools.go                    # Build tools imports
```

## API Usage

### Authentication Flow

The API uses **JWT Bearer token** authentication. Include the token in the `Authorization` header:

```http
Authorization: Bearer <your_jwt_token>
```

### User Management

#### Register a new user

```graphql
mutation {
  createUser(
    input: {
      email: "john.doe@example.com"
      fullname: "John Doe"
      password: "securePassword123"
    }
  )
}
```

**Response:**

```json
{
  "data": {
    "createUser": true
  }
}
}
```

#### Login

```graphql
mutation {
  login(input: { email: "john.doe@example.com", password: "securePassword123" })
}
```

**Response:**

```json
{
  "data": {
    "login": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
}
```

#### Refresh Token (Not implemented yet)

```graphql
mutation {
  refreshToken(input: { token: "your_refresh_token_here" })
}
```

### Notes Management

> **Note:** All note operations require authentication. Include the JWT token in the Authorization header.

#### Get all notes (Authenticated)

```graphql
query {
  notes {
    id
    title
    description
    text
  }
}
```

#### Get a specific note (Authenticated)

```graphql
query {
  note(id: "01931b2c-8b5a-7c90-9d4e-abc123def456") {
    id
    title
    description
    text
  }
}
```

#### Create a note (Authenticated)

```graphql
mutation {
  createNote(
    input: {
      title: "My First Note"
      description: "This is a sample note description"
      text: "Here's the main content of my note. It can be quite long and contain detailed information."
    }
  )
}
```

#### Update a note (Authenticated)

```graphql
mutation {
  updateNote(
    input: {
      noteid: "01931b2c-8b5a-7c90-9d4e-abc123def456"
      title: "Updated Note Title"
      description: "Updated description"
      text: "Updated note content with new information."
    }
  )
}
```

#### Delete a note (Authenticated)

```graphql
mutation {
  deleteNote(input: { id: "01931b2c-8b5a-7c90-9d4e-abc123def456" })
}
```

### Authentication Example Workflow

1. **Register a user:**

   ```bash
   curl -X POST http://localhost:8080/gq \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { createUser(input: {email: \"test@example.com\", fullname: \"Test User\", password: \"password123\"}) }"}'
   ```

2. **Login to get JWT token:**

   ```bash
   curl -X POST http://localhost:8080/gq \
     -H "Content-Type: application/json" \
     -d '{"query":"mutation { login(input: {email: \"test@example.com\", password: \"password123\"}) }"}'
   ```

3. **Use the token for authenticated requests:**

   ```bash
   curl -X POST http://localhost:8080/gq \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{"query":"query { notes { id title description text } }"}'
   ```

## Development Workflow

### 1. Schema-First Development

This project follows a **schema-first approach** for GraphQL development:

1. **Modify the GraphQL Schema:**
   Edit `graph/schema.graphqls` to add new types, queries, or mutations

   ```graphql
   # Example: Adding a new type
   type Comment {
     id: ID!
     text: String!
     noteId: ID!
     createdAt: String!
   }

   # Example: Adding a new query
   extend type Query {
     comments(noteId: ID!): [Comment!]!
   }
   ```

2. **Generate Code:**

   ```bash
   make gql-generate
   ```

   This creates resolver stubs for new operations.

3. **Implement Resolvers:**
   Add your business logic in `graph/schema.resolvers.go`

   ```go
   func (r *queryResolver) Comments(ctx context.Context, noteID string) ([]*model.Comment, error) {
       // Implement your logic here
       user := auth.ForContext(ctx)
       if user == nil {
           return nil, errors.New("unauthorized")
       }
       // ... business logic
   }
   ```

4. **Test in Playground:**
   Use the GraphQL playground at `http://localhost:8080/gq` to test your changes.

### 2. Database Changes

Follow this process for database schema changes:

1. **Create Migration Files:**

   Create up and down migration files in `internal/pkg/database/migrations/postgresql/`:

   ```bash
   # Example: Adding a comments table
   touch internal/pkg/database/migrations/postgresql/000003_create_comments_table.up.sql
   touch internal/pkg/database/migrations/postgresql/000003_create_comments_table.down.sql
   ```

2. **Write Migration SQL:**

   **Up migration** (`000003_create_comments_table.up.sql`):

   ```sql
   CREATE TABLE IF NOT EXISTS comments (
     id VARCHAR(36) PRIMARY KEY,
     text TEXT NOT NULL,
     note_id VARCHAR(36) NOT NULL,
     user_id VARCHAR(36) NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     FOREIGN KEY(note_id) REFERENCES notes(id) ON DELETE CASCADE,
     FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
   );
   ```

   **Down migration** (`000003_create_comments_table.down.sql`):

   ```sql
   DROP TABLE IF EXISTS comments;
   ```

3. **Apply Migrations:**

   Restart the application to automatically apply migrations, or run:

   ```bash
   go run server.go
   ```

### 3. Adding New Domains

To add a new business domain (e.g., comments, categories):

1. **Create Domain Package:**

   ```bash
   mkdir -p internal/domain/comments
   ```

2. **Define Repository Interface:**

   Create `internal/domain/comments/comments.go`:

   ```go
   package comments

   import (
       "context"
       "log/slog"
       "github.com/daffadon/graphy/graph/model"
       "github.com/daffadon/graphy/internal/infrastructure/database"
   )

   type CommentRepository interface {
       CreateComment(ctx context.Context, comment *model.Comment) error
       GetCommentsByNoteID(ctx context.Context, noteID, userID string) ([]*model.Comment, error)
       // ... other methods
   }

   type commentRepository struct {
       q database.Querier
       l *slog.Logger
   }

   func NewCommentRepository(q database.Querier, l *slog.Logger) CommentRepository {
       return &commentRepository{q: q, l: l}
   }
   ```

3. **Add to Dependency Injection:**

   Update `cmd/bootstrap.go`:

   ```go
   func BootstrapRun() *IBootstrap {
       // ... existing code
       cr := comments.NewCommentRepository(q, slog) // Add this line

       return &IBootstrap{
           G: &graph.Resolver{
               Ur: ur,
               Nr: nr,
               Cr: cr, // Add to resolver
               S:  slog,
           },
           // ... rest
       }
   }
   ```

4. **Update GraphQL Resolver:**

   Add the repository to `graph/resolver.go`:

   ```go
   type Resolver struct {
       Ur users.UserRepository
       Nr notes.NoteRepository
       Cr comments.CommentRepository // Add this line
       S  *slog.Logger
   }
   ```

### 4. Environment-Specific Configuration

The application supports different environments:

- **Development:** Uses `config.local.yaml` (default)
- **Production:** Uses `config.yaml` (set `ENV=production`)
- **Testing:** Uses `config.test.yaml` (set `ENV=test`)

Environment is determined by the `ENV` environment variable in `config/env/env.go`.

### 5. Code Organization Best Practices

- **Repository Pattern:** Each domain has its own repository interface
- **Dependency Injection:** All dependencies are injected through `bootstrap.go`
- **Error Handling:** Use structured logging with `slog` for consistent error reporting
- **Authentication:** Use the `auth.ForContext(ctx)` helper to get the authenticated user
- **Database Queries:** Use Squirrel query builder for type-safe SQL generation
- **Testing:** Write unit tests for repository and resolver functions

## Available Commands

### Makefile Commands

```bash
# Initialize GraphQL for first-time setup
make gql-init

# Generate GraphQL code after schema changes
make gql-generate
```

### Docker Commands

```bash
# Start only PostgreSQL
docker-compose up -d postgresql

# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View PostgreSQL logs
docker-compose logs postgresql

# Connect to PostgreSQL container
docker-compose exec postgresql psql -U postgres -d graphy_db

# Remove volumes (WARNING: This deletes all data)
docker-compose down -v
```

### Go Commands

```bash
# Run the application
go run server.go

# Build for current platform
go build -o tmp/main.exe server.go

# Build for Linux (from Windows)
GOOS=linux GOARCH=amd64 go build -o graphy server.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Download dependencies
go mod download

# Clean up dependencies
go mod tidy

# Update dependencies
go get -u ./...

# Verify dependencies
go mod verify

# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...
```

### Air Commands (Development)

```bash
# Start with hot reload (recommended for development)
air

# Start with custom config
air -c .air.toml

# Build only (without running)
air -build
```

## Configuration

The application uses a flexible YAML-based configuration system with environment-specific files:

### Configuration Files

- `config.example.yaml` - Template configuration with all available options
- `config.local.yaml` - Local development configuration (create this file)
- `config.yaml` - Production configuration
- `config.test.yaml` - Test environment configuration

### Configuration Structure

```yaml
database:
  sql:
    protocol: postgres # Database protocol
    host: localhost # Database host
    user: postgres # Database username
    password: your_password # Database password
    port: 5432 # Database port
    name: graphy_db # Database name
    sslmode: disable # SSL mode (disable/require/verify-full)

app:
  port: 8080 # Application port

server:
  cors:
    allow_origins: "*" # Allowed origins (* for all)
    allow_methods: "GET, POST, PUT, DELETE, OPTIONS, PATCH"
    allow_headers: "Content-Type, Authorization, X-Requested-With, X-CSRF-Token, Accept, Origin, Cache-Control, X-File-Name, X-File-Type, X-File-Size"
    expose_headers: "Content-Length, Content-Range"
    max_age: 86400 # Preflight cache duration in seconds
    allow_credential: true # Allow credentials in CORS requests
```

### Environment Variables

The application checks the `ENV` environment variable to determine which config file to load:

- `ENV=production` → loads `config.yaml`
- `ENV=test` → loads `config.test.yaml`
- Default (or `ENV=development`) → loads `config.local.yaml`

### Key Configuration Sections

#### Database Configuration (`database.sql`)

- **protocol**: PostgreSQL connection protocol
- **host**: Database server hostname or IP
- **user/password**: Database credentials
- **port**: Database server port (default: 5432)
- **name**: Database name
- **sslmode**: SSL connection mode for security

#### Application Configuration (`app`)

- **port**: HTTP server port (default: 8080)

#### CORS Configuration (`server.cors`)

- **allow_origins**: Comma-separated list of allowed origins
- **allow_methods**: HTTP methods allowed for CORS requests
- **allow_headers**: Headers allowed in CORS requests
- **expose_headers**: Headers exposed to the client
- **max_age**: How long browsers can cache preflight responses
- **allow_credential**: Whether to allow credentials (cookies, auth headers)

### JWT Configuration

JWT secret key is configured in the config file

```go
var secretKey = []byte(viper.GetString("jwt.s_key"))
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Error

**Problem:** `failed to create migrate instance` or connection refused errors

**Solutions:**

- Ensure PostgreSQL is running: `docker ps`
- Check database configuration in `config.yaml` for production or `config.local.yaml` for development
- Verify Docker container is healthy: `docker-compose logs postgresql`
- Test connection manually:
  ```bash
  docker-compose exec postgresql psql -U postgres -d graphy_db -c "SELECT 1;"
  ```
- Ensure database exists:
  ```bash
  docker-compose exec postgresql createdb -U postgres graphy_db
  ```

#### 2. GraphQL Generation Errors

**Problem:** `gqlgen` command fails or generates incorrect code

**Solutions:**

- Run `go mod tidy` to ensure dependencies are up to date
- Check `gqlgen.yml` configuration syntax
- Verify GraphQL schema syntax in `graph/schema.graphqls`
- Clear generated files and regenerate:
  ```bash
  rm -rf graph/generated.go graph/model/models_gen.go
  make gql-generate
  ```
- Ensure gqlgen version compatibility: `go list -m github.com/99designs/gqlgen`

#### 3. Migration Errors

**Problem:** Database migration fails on startup

**Solutions:**

- Check database permissions for the user
- Verify migration file SQL syntax
- Ensure database exists and is accessible
- Check migration file naming convention (must be sequential)
- Manually run migrations to see detailed error:
  ```bash
  go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
    -path internal/pkg/database/migrations/postgresql \
    -database "postgres://postgres:password@localhost:5432/graphy_db?sslmode=disable" up
  ```

#### 4. Port Already in Use

**Problem:** `bind: address already in use` error

**Solutions:**

- Change the port in `config.local.yaml`:
  ```yaml
  app:
    port: 8081 # Use different port
  ```
- Kill existing processes:

  ```bash
  # Windows
  netstat -ano | findstr :8080
  taskkill /PID <PID> /F

  # Linux/macOS
  lsof -ti:8080 | xargs kill -9
  ```

- Check if Air is already running: `ps aux | grep air`

#### 5. JWT Token Issues

**Problem:** Authentication fails or "Invalid token" errors

**Solutions:**

- Ensure token is included in Authorization header: `Authorization: Bearer <token>`
- Check token expiration (default: 2 hours)
- Verify JWT secret key consistency
- Test token generation:
  ```graphql
  mutation {
    login(input: { email: "test@example.com", password: "password123" })
  }
  ```

#### 6. CORS Issues

**Problem:** Browser blocks requests due to CORS policy

**Solutions:**

- Update CORS configuration in `config.local.yaml`:
  ```yaml
  server:
    cors:
      allow_origins: "http://localhost:3000,http://localhost:8080"
      allow_credentials: true
  ```
- Ensure preflight requests are handled correctly
- Check browser developer tools for specific CORS errors

#### 7. Hot Reload Not Working

**Problem:** Air doesn't detect file changes or restart

**Solutions:**

- Check `.air.toml` configuration
- Ensure file watchers aren't at OS limits:
  ```bash
  # Linux: Increase inotify watches
  echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf
  sudo sysctl -p
  ```
- Verify Air is watching correct directories
- Check `tmp/build-errors.log` for build issues

### Debug Mode

#### Enable Detailed Logging

Modify `config/logger/slog.go` to enable debug logging:

```go
func NewSlog() *slog.Logger {
    baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug, // Change from LevelInfo to LevelDebug
    })
    colorHandler := &ColorTextHandler{Handler: baseHandler}
    return slog.New(colorHandler)
}
```

#### GraphQL Query Logging

Enable query logging by modifying the server setup in `server.go`:

```go
srv.Use(extension.Logger{})
```

#### Database Query Logging

For PostgreSQL query logging, modify the connection in `config/database/postgresql.go`:

```go
config, err := pgxpool.ParseConfig(dsn)
if err != nil {
    logger.Error("Failed to parse database config")
}
config.ConnConfig.Tracer = &tracelog.TraceLog{
    Logger:   pgxlog.NewLogger(logger),
    LogLevel: tracelog.LogLevelDebug,
}
pool, err := pgxpool.NewWithConfig(context.Background(), config)
```

### Getting Help

1. **Check application logs** for detailed error messages
2. **Enable debug logging** to see internal operations
3. **Review GraphQL errors** in the playground
4. **Check database logs** for SQL-related issues
5. **Use Go's built-in race detector:** `go run -race server.go`
6. **Profile the application** during development with pprof

## Technology Stack

### Backend Technologies

- **[Go](https://golang.org/)** - Primary programming language (v1.24.6+)
- **[gqlgen](https://gqlgen.com/)** - GraphQL server library for Go
- **[Chi](https://github.com/go-chi/chi)** - Lightweight HTTP router
- **[PostgreSQL](https://postgresql.org/)** - Primary database
- **[pgx](https://github.com/jackc/pgx)** - PostgreSQL driver and toolkit
- **[Squirrel](https://github.com/Masterminds/squirrel)** - SQL query builder
- **[Viper](https://github.com/spf13/viper)** - Configuration management
- **[slog](https://pkg.go.dev/log/slog)** - Structured logging (Go 1.21+)

### Authentication & Security

- **[JWT](https://github.com/golang-jwt/jwt)** - JSON Web Token authentication
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Password hashing
- **CORS** - Cross-origin resource sharing support

### Database & Migrations

- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migration tool
- **[UUID v7](https://github.com/google/uuid)** - Time-ordered unique identifiers

### Development Tools

- **[Air](https://github.com/air-verse/air)** - Hot reload for development
- **[Docker](https://docker.com/)** - Containerization
- **[Docker Compose](https://docs.docker.com/compose/)** - Multi-container orchestration

### Architecture Patterns

- **Clean Architecture** - Separation of concerns with clear boundaries
- **Repository Pattern** - Data access abstraction
- **Dependency Injection** - Loose coupling between components
- **Domain-Driven Design** - Business logic organization

## Key Features Explained

### JWT Authentication Flow

1. **User Registration**: Password is hashed using bcrypt before storage
2. **Login**: Credentials are validated and JWT token is generated (2-hour expiry)
3. **Authorization**: Middleware extracts and validates JWT from Authorization header
4. **Context Injection**: User information is injected into GraphQL context
5. **Resolver Access**: Resolvers use `auth.ForContext(ctx)` to get authenticated user

### Database Design

#### Users Table

```sql
CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,          -- UUID v7
  email VARCHAR(255) NOT NULL UNIQUE, -- User email (unique)
  fullname VARCHAR(255) NOT NULL,     -- Full name
  password VARCHAR(255) NOT NULL,     -- bcrypt hashed password
  created_at TIMESTAMP DEFAULT NOW(), -- Creation timestamp
  updated_at TIMESTAMP DEFAULT NOW()  -- Last update timestamp
);
```

#### Notes Table

```sql
CREATE TABLE notes (
  id VARCHAR(36) PRIMARY KEY,            -- UUID v7
  title VARCHAR(255) NOT NULL,           -- Note title
  description TEXT NOT NULL,             -- Note description
  text TEXT NOT NULL,                    -- Note content
  user_id VARCHAR(36) NOT NULL,          -- Foreign key to users
  created_at TIMESTAMPTZ DEFAULT NOW(),  -- Creation timestamp with timezone
  updated_at TIMESTAMPTZ DEFAULT NOW(),  -- Last update timestamp with timezone
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### GraphQL Schema Design

The schema follows GraphQL best practices:

- **Mutations return booleans** for simple operations
- **Queries return typed objects** with all necessary fields
- **Input types** are used for complex operations
- **ID types** are used for unique identifiers
- **Non-null fields** are explicitly marked with `!`

### Security Features

1. **Password Security**: bcrypt hashing with default cost factor
2. **JWT Security**: HS256 signing algorithm with expiration
3. **Authorization**: All note operations require valid JWT
4. **Data Isolation**: Users can only access their own notes
5. **CORS Protection**: Configurable CORS policies
6. **SQL Injection Prevention**: Parameterized queries via Squirrel

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [gqlgen](https://gqlgen.com/) for the excellent GraphQL library
- [Chi](https://github.com/go-chi/chi) for the lightweight HTTP router
- [pgx](https://github.com/jackc/pgx) for the robust PostgreSQL driver
- [Air](https://github.com/air-verse/air) for hot reload functionality
- The Go community for excellent tooling and libraries

---

**Happy coding! 🚀**
