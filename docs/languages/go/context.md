# Context in Go

## Overview
This document explains the usage of Go's `context` package in the Kasho application. Contexts are used to carry deadlines, cancellation signals, and request-scoped values across API boundaries and between processes.

## Core Concepts

### Basic Context Creation
```go
// Root context - use as the base for all other contexts
ctx := context.Background()

// Placeholder context - use when unsure which context to use
ctx := context.TODO()

// Context with cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // Important: Always call cancel to prevent leaks

// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Context with deadline
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// Context with values
ctx := context.WithValue(context.Background(), "userID", "123")
```

## Key Features

### 1. Cancellation
Contexts allow you to cancel operations when they're no longer needed. This is particularly useful for:
- HTTP requests that time out
- Database queries that should be canceled
- Long-running operations that need to be stopped

### 2. Timeouts
Timeouts help prevent operations from running indefinitely. They're crucial for:
- Database operations
- External API calls
- Resource-intensive operations

### 3. Values
Contexts can carry request-scoped values through the call chain. Useful for:
- Request IDs for tracing
- User authentication information
- Operation metadata

## Best Practices

### Function Signatures
```go
// Good: Context as first parameter
func (q *Queries) GetUser(ctx context.Context, id string) (User, error)

// Bad: Context as struct field
type Queries struct {
    ctx context.Context  // Don't do this!
}
```

### Error Handling
```go
func (h *Handler) ProcessTransaction(ctx context.Context, tx Transaction) error {
    // Check context before starting work
    if ctx.Err() != nil {
        return ctx.Err()
    }
    
    // Check context after long operations
    if err := h.validateTransaction(ctx, tx); err != nil {
        if ctx.Err() != nil {
            return ctx.Err()  // Return context error if context is done
        }
        return err
    }
    
    return nil
}
```

### Resource Cleanup
```go
func (h *Handler) LongRunningOperation(ctx context.Context) error {
    // Create a child context with timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()  // Always defer cancel
    
    // Use the context
    done := make(chan error, 1)
    go func() {
        // Do work...
        done <- nil
    }()
    
    // Wait for either context cancellation or operation completion
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

## Common Pitfalls to Avoid

1. **Don't store contexts in structs**
   - Contexts should be passed as function parameters
   - They're designed for request-scoped data

2. **Always call `cancel()`**
   - When using `WithCancel`, `WithTimeout`, or `WithDeadline`
   - Use `defer cancel()` to ensure cleanup

3. **Don't pass `nil` as a context**
   - Use `context.Background()` or `context.TODO()`
   - `Background()` is the root context
   - `TODO()` indicates a placeholder

4. **Don't use context values for optional parameters**
   - Use function parameters instead
   - Context values should be for request-scoped data

5. **Don't use context values for data that should be in function parameters**
   - Keep function signatures explicit
   - Use context for cross-cutting concerns

## Usage in Kasho

### Database Operations
Contexts are used in database operations to:
- Set timeouts for queries
- Cancel long-running operations
- Pass transaction information

### API Handlers
Contexts in HTTP handlers help with:
- Request timeouts
- Client disconnection handling
- Request tracing
- Authentication information

### Background Jobs
Contexts in background jobs enable:
- Graceful shutdown
- Job cancellation
- Timeout management

## Examples

### Database Transactions with Context
This example shows how to use context in a financial transaction, including timeouts, cancellation, and proper error handling:

```go
// Example of a money transfer with context
func (h *Handler) TransferMoney(ctx context.Context, from, to uuid.UUID, amount decimal.Decimal) error {
    // Create a child context with timeout for the entire transaction
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Start transaction
    tx, err := h.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable, // Highest isolation level for financial transactions
    })
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback() // Will be ignored if tx.Commit() is called

    // Get account balances with context
    fromAccount, err := h.queries.WithTx(tx).GetAccount(ctx, from)
    if err != nil {
        return fmt.Errorf("failed to get source account: %w", err)
    }

    toAccount, err := h.queries.WithTx(tx).GetAccount(ctx, to)
    if err != nil {
        return fmt.Errorf("failed to get destination account: %w", err)
    }

    // Check if context is done before proceeding
    if ctx.Err() != nil {
        return ctx.Err()
    }

    // Perform transfer
    if err := h.queries.WithTx(tx).DeductBalance(ctx, from, amount); err != nil {
        return fmt.Errorf("failed to deduct balance: %w", err)
    }

    if err := h.queries.WithTx(tx).AddBalance(ctx, to, amount); err != nil {
        return fmt.Errorf("failed to add balance: %w", err)
    }

    // Record transaction
    if err := h.queries.WithTx(tx).RecordTransaction(ctx, db.RecordTransactionParams{
        FromAccountID: from,
        ToAccountID:   to,
        Amount:        amount,
        Type:          "transfer",
    }); err != nil {
        return fmt.Errorf("failed to record transaction: %w", err)
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}
```

Key points in this example:
- Uses `WithTimeout` to set a deadline for the entire transaction
- Properly handles transaction rollback with `defer`
- Checks context status before critical operations
- Uses proper error wrapping with `fmt.Errorf`
- Implements serializable isolation level for financial safety

### API Handlers with Context
This example demonstrates context usage in an HTTP handler, including request timeouts and proper error handling:

```go
// Example of an API handler with context
func (h *Handler) CreateAccount(c *gin.Context) {
    // Get the request context
    ctx := c.Request.Context()

    // Add timeout for the entire operation
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // Parse request body
    var req struct {
        AccountType string `json:"account_type" binding:"required"`
        Currency    string `json:"currency" binding:"required,len=3"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get user ID from context (set by auth middleware)
    userID, ok := c.Get("userID")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    // Create account with context
    account, err := h.queries.CreateAccount(ctx, db.CreateAccountParams{
        UserID:      userID.(uuid.UUID),
        AccountType: req.AccountType,
        Currency:    req.Currency,
    })
    if err != nil {
        // Check if error is due to context cancellation
        if ctx.Err() == context.DeadlineExceeded {
            c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
            return
        }
        if ctx.Err() == context.Canceled {
            c.JSON(http.StatusServiceUnavailable, gin.H{"error": "request canceled"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
        return
    }

    // Add request ID to response headers
    c.Header("X-Request-ID", ctx.Value("requestID").(string))

    c.JSON(http.StatusCreated, account)
}
```

Key points in this example:
- Uses request context from Gin
- Implements operation timeout
- Handles different types of context errors
- Properly propagates request ID
- Uses context for authentication

### Middleware with Context
This example shows how to use context in middleware for authentication and request tracking:

```go
// Example of authentication middleware with context
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the request context
        ctx := c.Request.Context()

        // Add request ID to context
        requestID := uuid.New().String()
        ctx = context.WithValue(ctx, "requestID", requestID)
        c.Request = c.Request.WithContext(ctx)

        // Get token from header
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            c.Abort()
            return
        }

        // Validate token with timeout
        ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
        defer cancel()

        claims, err := h.validateToken(ctx, token)
        if err != nil {
            if ctx.Err() == context.DeadlineExceeded {
                c.JSON(http.StatusGatewayTimeout, gin.H{"error": "token validation timed out"})
            } else {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            }
            c.Abort()
            return
        }

        // Add user ID to context
        ctx = context.WithValue(ctx, "userID", claims.UserID)
        c.Request = c.Request.WithContext(ctx)

        c.Next()
    }
}
```

Key points in this example:
- Adds request ID for tracing
- Implements token validation with timeout
- Properly propagates context through the request
- Handles different types of errors
- Uses context for user authentication

### Background Jobs with Context
This example demonstrates how to use context in background jobs for graceful shutdown and timeout management:

```go
// Example of a background job processor with context
type TransactionProcessor struct {
    db      *sql.DB
    queries *db.Queries
    logger  *zap.Logger
}

func (p *TransactionProcessor) Start(ctx context.Context) error {
    // Create a ticker for periodic processing
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            // Context was canceled, perform graceful shutdown
            p.logger.Info("shutting down transaction processor",
                zap.String("reason", ctx.Err().Error()))
            return ctx.Err()
        case <-ticker.C:
            // Process pending transactions
            if err := p.processPendingTransactions(ctx); err != nil {
                if ctx.Err() != nil {
                    // Context was canceled during processing
                    return ctx.Err()
                }
                p.logger.Error("failed to process transactions",
                    zap.Error(err))
                // Continue processing despite error
                continue
            }
        }
    }
}

func (p *TransactionProcessor) processPendingTransactions(ctx context.Context) error {
    // Create a child context with timeout for this batch
    ctx, cancel := context.WithTimeout(ctx, 55*time.Second)
    defer cancel()

    // Get pending transactions
    transactions, err := p.queries.GetPendingTransactions(ctx)
    if err != nil {
        return fmt.Errorf("failed to get pending transactions: %w", err)
    }

    // Process each transaction
    for _, tx := range transactions {
        // Check if context is done before processing each transaction
        if ctx.Err() != nil {
            return ctx.Err()
        }

        // Create a child context for individual transaction processing
        txCtx, txCancel := context.WithTimeout(ctx, 5*time.Second)
        if err := p.processTransaction(txCtx, tx); err != nil {
            txCancel()
            p.logger.Error("failed to process transaction",
                zap.String("transaction_id", tx.ID.String()),
                zap.Error(err))
            continue
        }
        txCancel()
    }

    return nil
}

func (p *TransactionProcessor) processTransaction(ctx context.Context, tx db.Transaction) error {
    // Start transaction
    dbTx, err := p.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer dbTx.Rollback()

    // Process the transaction with context
    if err := p.queries.WithTx(dbTx).ProcessTransaction(ctx, tx.ID); err != nil {
        return fmt.Errorf("failed to process transaction: %w", err)
    }

    // Commit the transaction
    if err := dbTx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

// Usage in main function
func main() {
    // Create root context
    ctx := context.Background()

    // Create a context that can be canceled
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // Handle OS signals for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel() // Cancel context on signal
    }()

    // Start the processor
    processor := NewTransactionProcessor(db, queries, logger)
    if err := processor.Start(ctx); err != nil {
        if err != context.Canceled {
            logger.Fatal("processor failed", zap.Error(err))
        }
        logger.Info("processor stopped gracefully")
    }
}
```

Key points in this example:
- Uses context for graceful shutdown
- Implements periodic processing with timeouts
- Handles OS signals for clean termination
- Uses child contexts for different processing stages
- Properly propagates context through the call chain
- Implements proper error handling and logging

### Request Tracing with Context
This example shows how to implement request tracing using context:

```go
// Request tracing middleware
func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get or generate request ID
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }

        // Create a new context with tracing information
        ctx := c.Request.Context()
        ctx = context.WithValue(ctx, "requestID", requestID)
        ctx = context.WithValue(ctx, "startTime", time.Now())
        
        // Add tracing headers
        c.Header("X-Request-ID", requestID)
        c.Request = c.Request.WithContext(ctx)

        // Create a span for this request
        span, ctx := opentracing.StartSpanFromContext(ctx, c.Request.URL.Path)
        defer span.Finish()

        // Add request details to span
        span.SetTag("http.method", c.Request.Method)
        span.SetTag("http.url", c.Request.URL.String())
        span.SetTag("http.request_id", requestID)

        // Process the request
        c.Next()

        // Add response details to span
        span.SetTag("http.status_code", c.Writer.Status())
        span.SetTag("http.duration_ms", time.Since(ctx.Value("startTime").(time.Time)).Milliseconds())

        // Log request details
        logger.Info("request completed",
            zap.String("request_id", requestID),
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", time.Since(ctx.Value("startTime").(time.Time))),
        )
    }
}

// Example of using tracing in a service
type AccountService struct {
    tracer opentracing.Tracer
    db     *sql.DB
}

func (s *AccountService) GetAccountBalance(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error) {
    // Start a new span for this operation
    span, ctx := opentracing.StartSpanFromContext(ctx, "GetAccountBalance")
    defer span.Finish()

    // Add operation details to span
    span.SetTag("account_id", accountID.String())
    span.SetTag("operation", "get_balance")

    // Get request ID from context
    if requestID, ok := ctx.Value("requestID").(string); ok {
        span.SetTag("request_id", requestID)
    }

    // Perform database operation
    balance, err := s.db.GetAccountBalance(ctx, accountID)
    if err != nil {
        // Log error to span
        span.SetTag("error", true)
        span.LogFields(log.String("error", err.Error()))
        return decimal.Zero, fmt.Errorf("failed to get account balance: %w", err)
    }

    // Log success to span
    span.SetTag("balance", balance.String())
    return balance, nil
}

// Example of using tracing in a database operation
func (db *DB) GetAccountBalance(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error) {
    // Start a new span for database operation
    span, ctx := opentracing.StartSpanFromContext(ctx, "DB.GetAccountBalance")
    defer span.Finish()

    // Add database operation details to span
    span.SetTag("db.operation", "SELECT")
    span.SetTag("db.table", "accounts")
    span.SetTag("db.account_id", accountID.String())

    // Perform the query
    var balance decimal.Decimal
    err := db.QueryRowContext(ctx, 
        "SELECT balance FROM accounts WHERE id = $1",
        accountID,
    ).Scan(&balance)

    if err != nil {
        span.SetTag("error", true)
        span.LogFields(log.String("error", err.Error()))
        return decimal.Zero, err
    }

    span.SetTag("db.result", balance.String())
    return balance, nil
}
```

Key points in this example:
- Implements distributed tracing using OpenTracing
- Shows how to propagate request IDs through context
- Demonstrates creating and managing spans for different operations
- Includes proper error logging and span tagging
- Shows how to track request duration and status
- Implements middleware for request tracing
- Shows how to use tracing in services and database operations

### Rate Limiting with Context
This example demonstrates how to implement rate limiting using context for both per-user and global rate limits:

```go
// Rate limiter implementation using token bucket algorithm
type RateLimiter struct {
    mu       sync.Mutex
    limiters map[string]*rate.Limiter
    global   *rate.Limiter
    logger   *zap.Logger
}

func NewRateLimiter(globalRPS float64, perUserRPS float64) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        global:   rate.NewLimiter(rate.Limit(globalRPS), int(globalRPS)),
        logger:   zap.L().Named("rate_limiter"),
    }
}

// Rate limiting middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        
        // Get user ID from context (set by auth middleware)
        userID, ok := ctx.Value("userID").(string)
        if !ok {
            userID = "anonymous"
        }

        // Check global rate limit
        if !rl.global.Allow() {
            rl.logger.Warn("global rate limit exceeded",
                zap.String("user_id", userID),
                zap.String("path", c.Request.URL.Path))
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
                "retry_after": time.Until(rl.global.Reserve().DelayFrom(time.Now())).Seconds(),
            })
            c.Abort()
            return
        }

        // Get or create user-specific limiter
        limiter := rl.getUserLimiter(userID)

        // Try to acquire token with context
        if err := limiter.Wait(ctx); err != nil {
            if ctx.Err() == context.DeadlineExceeded {
                c.JSON(http.StatusGatewayTimeout, gin.H{"error": "rate limit wait timeout"})
            } else if ctx.Err() == context.Canceled {
                c.JSON(http.StatusServiceUnavailable, gin.H{"error": "request canceled"})
            } else {
                c.JSON(http.StatusTooManyRequests, gin.H{
                    "error": "rate limit exceeded",
                    "retry_after": time.Until(limiter.Reserve().DelayFrom(time.Now())).Seconds(),
                })
            }
            c.Abort()
            return
        }

        // Add rate limit headers
        c.Header("X-RateLimit-Limit", fmt.Sprintf("%.0f", limiter.Limit()))
        c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Tokens()))
        c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Until(limiter.Reserve().DelayFrom(time.Now()))).Unix()))

        c.Next()
    }
}

func (rl *RateLimiter) getUserLimiter(userID string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    limiter, exists := rl.limiters[userID]
    if !exists {
        limiter = rate.NewLimiter(rate.Limit(10), 10) // 10 requests per second, burst of 10
        rl.limiters[userID] = limiter
    }
    return limiter
}

// Example of using rate limiter in a service
type PaymentService struct {
    limiter *RateLimiter
    db      *sql.DB
}

func (s *PaymentService) ProcessPayment(ctx context.Context, payment Payment) error {
    // Create a child context with timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // Get user ID from context
    userID, ok := ctx.Value("userID").(string)
    if !ok {
        return fmt.Errorf("user ID not found in context")
    }

    // Check rate limit with context
    limiter := s.limiter.getUserLimiter(userID)
    if err := limiter.Wait(ctx); err != nil {
        return fmt.Errorf("rate limit exceeded: %w", err)
    }

    // Process payment
    // ... payment processing logic ...
    return nil
}
```

Key points in this example:
- Implements both global and per-user rate limiting using token bucket algorithm
- Uses context for timeout and cancellation in rate limiting
- Includes proper rate limit headers and retry-after information
- Demonstrates rate limiting middleware for HTTP requests
- Shows how to use rate limiting in services
- Implements proper error handling and logging
- Provides graceful handling of rate limit exceeded scenarios
- Supports both anonymous and authenticated users

### Caching with Context
This example shows how to implement caching with context for both in-memory and distributed caching:

```go
// Cache interface
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

// Redis cache implementation
type RedisCache struct {
    client *redis.Client
    logger *zap.Logger
}

func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
    // Create a child context with timeout
    ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
    defer cancel()

    // Get value from Redis with context
    val, err := c.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, cache.ErrNotFound
        }
        return nil, fmt.Errorf("redis get error: %w", err)
    }

    return val, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
    defer cancel()

    return c.client.Set(ctx, key, value, ttl).Err()
}

// In-memory cache with context
type MemoryCache struct {
    mu     sync.RWMutex
    items  map[string]cacheItem
    logger *zap.Logger
}

type cacheItem struct {
    value      interface{}
    expiresAt  time.Time
}

func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
    // Check if context is done
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }

    c.mu.RLock()
    defer c.mu.RUnlock()

    item, exists := c.items[key]
    if !exists {
        return nil, cache.ErrNotFound
    }

    if time.Now().After(item.expiresAt) {
        go c.Delete(context.Background(), key) // Delete expired item
        return nil, cache.ErrNotFound
    }

    return item.value, nil
}

// Cache middleware
func CacheMiddleware(cache Cache, ttl time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip caching for non-GET requests
        if c.Request.Method != http.MethodGet {
            c.Next()
            return
        }

        // Generate cache key
        key := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.String())

        // Try to get from cache
        if val, err := cache.Get(c.Request.Context(), key); err == nil {
            c.JSON(http.StatusOK, val)
            c.Abort()
            return
        }

        // Create a custom response writer to capture the response
        writer := &responseWriter{
            ResponseWriter: c.Writer,
            body:          &bytes.Buffer{},
        }
        c.Writer = writer

        // Process the request
        c.Next()

        // Cache the response if it was successful
        if c.Writer.Status() == http.StatusOK {
            var response interface{}
            if err := json.Unmarshal(writer.body.Bytes(), &response); err == nil {
                go cache.Set(context.Background(), key, response, ttl)
            }
        }
    }
}

// Example of using cache in a service
type AccountService struct {
    cache Cache
    db    *sql.DB
}

func (s *AccountService) GetAccount(ctx context.Context, accountID uuid.UUID) (*Account, error) {
    // Try to get from cache first
    cacheKey := fmt.Sprintf("account:%s", accountID.String())
    if val, err := s.cache.Get(ctx, cacheKey); err == nil {
        if account, ok := val.(*Account); ok {
            return account, nil
        }
    }

    // If not in cache, get from database
    account, err := s.db.GetAccount(ctx, accountID)
    if err != nil {
        return nil, err
    }

    // Cache the result
    go s.cache.Set(context.Background(), cacheKey, account, 5*time.Minute)

    return account, nil
}

// Example of using cache with background refresh
func (s *AccountService) GetAccountWithRefresh(ctx context.Context, accountID uuid.UUID) (*Account, error) {
    cacheKey := fmt.Sprintf("account:%s", accountID.String())
    
    // Try to get from cache
    if val, err := s.cache.Get(ctx, cacheKey); err == nil {
        if account, ok := val.(*Account); ok {
            // Start background refresh if cache is about to expire
            if time.Until(account.UpdatedAt.Add(5 * time.Minute)) < time.Minute {
                go s.refreshAccountCache(context.Background(), accountID)
            }
            return account, nil
        }
    }

    // If not in cache, get from database
    return s.refreshAccountCache(ctx, accountID)
}

func (s *AccountService) refreshAccountCache(ctx context.Context, accountID uuid.UUID) (*Account, error) {
    account, err := s.db.GetAccount(ctx, accountID)
    if err != nil {
        return nil, err
    }

    // Cache the result
    go s.cache.Set(context.Background(), 
        fmt.Sprintf("account:%s", accountID.String()),
        account,
        5*time.Minute,
    )

    return account, nil
}
```

Key points in this example:
- Implements both in-memory and Redis caching
- Uses context for timeouts and cancellation
- Implements cache middleware for HTTP responses
- Shows background cache refresh
- Handles cache misses and errors
- Implements proper TTL management
- Uses context for request-scoped operations
- Includes proper error handling and logging

## Additional Resources

- [Go Context Package Documentation](https://pkg.go.dev/context)
- [Go Blog: Context](https://go.dev/blog/context)
- [Go Context Best Practices](https://pkg.go.dev/context#pkg-overview) 