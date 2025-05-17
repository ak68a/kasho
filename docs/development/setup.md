# Development Guide

This guide will help you get started with development on the Kasho project.

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Node.js 18 or later
- Yarn package manager
- Git
- Make (optional but recommended)

## Initial Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/kasho.git
   cd kasho
   ```

2. **Backend Setup**
   ```bash
   cd backend
   
   # Install Go dependencies
   go mod tidy
   
   # Start the database
   make p_up
   
   # Create the database
   make db_up
   
   # Run migrations
   make m_up
   
   # Generate database code
   make sqlc
   
   # Run tests
   go test ./...
   ```

3. **Frontend Setup**
   ```bash
   cd frontend
   
   # Install dependencies
   yarn install
   
   # Start development server
   yarn dev
   ```

## Development Workflow

### Backend Development

1. **Database Changes**
   - Create new migrations: `make c_m name=your_migration_name`
   - Apply migrations: `make m_up`
   - Rollback migrations: `make m_down`
   - Generate code: `make sqlc`

2. **Testing**
   - Run all tests: `go test ./...`
   - Run specific test: `go test ./db/tests -run TestName`
   - Run with coverage: `go test ./... -cover`

3. **Code Style**
   - Use `gofmt` for formatting
   - Follow Go best practices
   - Write tests for new features
   - Document public functions

### Frontend Development

1. **Development Server**
   - Start: `yarn dev`
   - Build: `yarn build`
   - Production: `yarn start`

2. **Code Style**
   - Use TypeScript
   - Follow Next.js best practices
   - Write component tests
   - Use proper typing

## Common Tasks

### Adding a New Database Table

1. Create migration:
   ```bash
   make c_m name=create_new_table
   ```

2. Edit the migration file in `backend/db/migrations/`

3. Run migration:
   ```bash
   make m_up
   ```

4. Add queries in `backend/db/queries/`

5. Generate code:
   ```bash
   make sqlc
   ```

### Adding a New API Endpoint

1. Create new handler in backend
2. Add tests
3. Update API documentation
4. Add frontend integration

### Database Operations

- Start database: `make p_up`
- Stop database: `make p_down`
- Create database: `make db_up`
- Drop database: `make db_down`

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check if Docker is running
   - Verify database is up: `make p_up`
   - Check connection string

2. **Migration Issues**
   - Check migration files
   - Try rolling back: `make m_down`
   - Reapply: `make m_up`

3. **Frontend Issues**
   - Clear node_modules: `rm -rf node_modules`
   - Reinstall: `yarn install`
   - Clear Next.js cache: `rm -rf .next`

## Best Practices

### Code Organization

- Keep related code together
- Use clear package structure
- Follow Go project layout
- Maintain separation of concerns

### Testing

- Write tests for new features
- Maintain test coverage
- Use table-driven tests
- Mock external dependencies

### Git Workflow

- Use feature branches
- Write clear commit messages
- Keep commits focused
- Review code before merging

### Documentation

- Document public APIs
- Update README files
- Keep documentation current
- Add comments for complex logic

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Next.js Documentation](https://nextjs.org/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)

## Getting Help

- Check existing documentation
- Review code comments
- Ask team members
- Open an issue for bugs
- Create a discussion for questions

## Code Examples

### Backend Examples

#### Creating a New Database Migration
```sql
-- backend/db/migrations/000002_create_accounts.up.sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    account_type VARCHAR(50) NOT NULL,
    balance DECIMAL(19,4) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_created_at ON accounts(created_at);
```

```sql
-- backend/db/migrations/000002_create_accounts.down.sql
DROP TABLE IF EXISTS accounts;
```

#### Adding a New Query
```sql
-- backend/db/queries/accounts.sql
-- name: CreateAccount :one
INSERT INTO accounts (
    user_id,
    account_type,
    currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserAccounts :many
SELECT * FROM accounts
WHERE user_id = $1
ORDER BY created_at DESC;
```

#### Implementing an API Handler
```go
// backend/api/accounts.go
func (h *Handler) CreateAccount(c *gin.Context) {
    var req struct {
        AccountType string `json:"account_type" binding:"required"`
        Currency    string `json:"currency" binding:"required,len=3"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID := getUserIDFromContext(c)
    
    account, err := h.queries.CreateAccount(c.Request.Context(), db.CreateAccountParams{
        UserID:      userID,
        AccountType: req.AccountType,
        Currency:    req.Currency,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
        return
    }
    
    c.JSON(http.StatusCreated, account)
}
```

### Frontend Examples

#### Creating a New API Client Function
```typescript
// frontend/lib/api/accounts.ts
export async function createAccount(data: {
  accountType: string;
  currency: string;
}): Promise<Account> {
  const response = await fetch('/api/v1/accounts', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    throw new Error('Failed to create account');
  }
  
  return response.json();
}
```

#### Creating a New React Component
```typescript
// frontend/components/AccountForm.tsx
import { useState } from 'react';
import { createAccount } from '@/lib/api/accounts';

export function AccountForm() {
  const [accountType, setAccountType] = useState('');
  const [currency, setCurrency] = useState('');
  const [error, setError] = useState<string | null>(null);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createAccount({ accountType, currency });
      // Handle success
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create account');
    }
  };
  
  return (
    <form onSubmit={handleSubmit}>
      {/* Form fields */}
    </form>
  );
}
```

### Testing Examples

#### Backend Test
```go
// backend/db/tests/accounts_test.go
func TestCreateAccount(t *testing.T) {
    ctx := context.Background()
    db := setupTestDB(t)
    queries := db.NewQueries()
    
    user := createTestUser(t, queries)
    
    account, err := queries.CreateAccount(ctx, db.CreateAccountParams{
        UserID:      user.ID,
        AccountType: "checking",
        Currency:    "USD",
    })
    
    require.NoError(t, err)
    assert.Equal(t, user.ID, account.UserID)
    assert.Equal(t, "checking", account.AccountType)
    assert.Equal(t, "USD", account.Currency)
    assert.Equal(t, decimal.NewFromInt(0), account.Balance)
}
```

#### Frontend Test
```typescript
// frontend/components/__tests__/AccountForm.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { AccountForm } from '../AccountForm';

describe('AccountForm', () => {
  it('creates an account successfully', async () => {
    render(<AccountForm />);
    
    fireEvent.change(screen.getByLabelText('Account Type'), {
      target: { value: 'checking' },
    });
    
    fireEvent.change(screen.getByLabelText('Currency'), {
      target: { value: 'USD' },
    });
    
    fireEvent.click(screen.getByText('Create Account'));
    
    await screen.findByText('Account created successfully');
  });
}); 