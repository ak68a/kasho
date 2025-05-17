# Kasho

A modern fintech application built with Go and React.

## Prerequisites

### Backend
- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL 15 (handled via Docker)
- sqlc (for database code generation)
- CompileDaemon (for hot reload during development)
- Key Go packages (automatically installed via `go mod tidy`):
  - Gin (web framework)
  - Viper (configuration management)
  - lib/pq (PostgreSQL driver)
- make (optional, for using make commands)

### Frontend
- Node.js 18 or later
- Yarn (package manager)

## Getting Started

### Backend Setup

1. Install Go dependencies:
```bash
cd backend
go mod tidy
```

2. Start the database (choose one):
```bash
# Using make
make p_up

# OR using docker-compose directly
docker-compose up -d
```

3. Create the database (if not exists):
```bash
make db_up
```

4. Run database migrations to create tables:
```bash
make m_up
```

5. Generate database code (only needed if you modify SQL queries or database schema):
```bash
# Using make
make sqlc

# OR manually
# Install sqlc if you haven't already
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# Generate the code
sqlc generate
```

Note: The generated code is committed to the repository, so you only need to run this if you:
- Modify SQL queries in `backend/db/queries/`
- Add new migrations that change the database schema
- Pull changes that include modified SQL queries or schema changes

6. Run the tests:
```bash
go test ./...
```

7. Run the API server:
```bash
# Without hot reload using CompileDaemon
go run main.go

# Using make
make start

# OR manually
go get github.com/githubnemo/CompileDaemon@latest
go install github.com/githubnemo/CompileDaemon@latest
CompileDaemon -command="./backend"

# Troubleshooting
#---
# zsh: command not found: CompileDaemon -> run
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc

# Or manually add to your ~/.zshrc file
export PATH="$PATH:$(go env GOPATH)/bin"

#---
# If CompileDaemon runs but hands on "Build ok." -> run
CompileDaemon -command="./backend" -build="go build -o backend main.go"

```

### Frontend Setup

1. Install dependencies:
```bash
cd frontend
yarn install
```

2. Start the Next.js development server:
```bash
yarn dev
```

The frontend will be available at `http://localhost:3000`

## Project Structure

```
.
├── backend/
├── ├── api/             
│   ├── db/
│   │   ├── migrations/   # Database migrations
│   │   ├── queries/      # SQL queries
│   │   ├── sqlc/         # Generated database code
│   │   └── tests/        # Database tests
│   └── utils/            # Utility functions
└── frontend/             # React frontend
```

## Database

The project uses PostgreSQL with the following default credentials:
- Host: localhost
- Port: 5432
- User: root
- Password: root
- Database: kasho_db

## Development

### Using Make Commands

The project includes several make commands for common tasks:

```bash
# Database Management
make p_up        # Start PostgreSQL server
make p_down      # Stop PostgreSQL server
make db_up       # Create database
make db_down     # Drop database

# Migrations
make c_m name=create_users    # Create a new migration
make m_up                     # Run migrations up
make m_down                   # Roll back migrations

# Code Generation
make sqlc                     # Generate database code

# API Server
make start                    # Run API server w/ hot reload
```

### Database Migrations

Database migrations are managed through SQL files in `backend/db/migrations/`. After adding new migrations:

1. Apply the migrations using make:
```bash
make m_up
```

2. Generate the code:
```bash
make sqlc
```

### Testing

Run the test suite:
```bash
# Using make
make test

# OR manually
cd backend
go test ./...

# Test w/ verbose output and coverage report
go test -v -cover ./...
```

