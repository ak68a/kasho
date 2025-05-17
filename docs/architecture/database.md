# Database Documentation

## Overview
This document describes the database schema, relationships, and management procedures for the Kasho application.

## Database Technology
- PostgreSQL 15
- Managed through Docker
- Migrations using golang-migrate
- Code generation using SQLC

## Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Accounts Table
```sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    account_type VARCHAR(50) NOT NULL,
    balance DECIMAL(19,4) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    account_id UUID REFERENCES accounts(id),
    transaction_type VARCHAR(50) NOT NULL,
    amount DECIMAL(19,4) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Relationships
- One-to-Many: Users to Accounts
- One-to-Many: Accounts to Transactions
- Many-to-One: Transactions to Accounts

## Indexes
- Primary keys
- Foreign keys
- Performance indexes
- Unique constraints

## Data Types
- UUID for IDs
- DECIMAL(19,4) for monetary values
- TIMESTAMP WITH TIME ZONE for dates
- VARCHAR for strings
- ENUM for status fields

## Migrations
- Location: `backend/db/migrations/`
- Naming: `{version}_{description}.up.sql`
- Rollback: `{version}_{description}.down.sql`

## Backup and Recovery
- Backup procedures
- Recovery procedures
- Point-in-time recovery
- Data retention policies

## Performance
- Query optimization
- Index usage
- Connection pooling
- Caching strategy

## Security
- Access control
- Encryption at rest
- Audit logging
- Data masking

## Monitoring
- Query performance
- Connection usage
- Disk space
- Backup status

## Maintenance
- Vacuum procedures
- Index maintenance
- Statistics updates
- Connection management 