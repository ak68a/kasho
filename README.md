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
  - gin-contrib/cors (CORS middleware)
- make (optional, for using make commands)

### Frontend
- Node.js 18 or later
- Yarn (package manager)
- Key frontend packages (automatically installed via `yarn install`):
  - Next.js (React framework)
  - Axios (HTTP client)
  - React Query (data fetching)
  - Tailwind CSS (styling)

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
# Using make
make start

# Without hot reload using CompileDaemon
go run main.go

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

### API Integration

#### CORS Configuration
The backend API is configured with CORS to allow requests from the frontend development server. The CORS middleware is set up in the backend to:
- Allow requests from `http://localhost:3000` (frontend development server)
- Support common HTTP methods (GET, POST, PUT, DELETE)
- Allow necessary headers for authentication and content negotiation

If you need to modify CORS settings, update the configuration in `backend/api/middleware/cors.go`.

#### Axios Setup
The frontend uses Axios for making HTTP requests to the backend API. The Axios instance is configured in `frontend/lib/axios.ts` with:
- Base URL configuration for different environments
- Default headers for content type and authentication
- Request/response interceptors for error handling
- Automatic token management

API routes are managed through URL constants in `frontend/utils/network.ts` to maintain consistency and avoid hardcoded strings:

```typescript
// Example from frontend/utils/network.ts
export const authUrl = {
  login: '/auth/login',
  register: '/auth/register',
  users: '/users',
  // ... other auth-related endpoints
} as const;
```

Example usage in frontend components:
```typescript
import axios from '@/lib/axios';

// Making API requests using URL constants
const response = await axios.get(authUrl.users);
const data = await axios.post(authUrl.login, credentials);
```

### Checking Ports

```bash
lsof -i :3000 # or any other port you might be using
```