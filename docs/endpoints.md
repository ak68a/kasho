# API Documentation

## Overview
This document outlines the API endpoints, request/response formats, and authentication requirements for the Kasho application.

## Base URL
```
http://localhost:3000/  # Development
https://api.kasho.com/  # Production
```

## Authentication
- JWT-based authentication
- Token refresh mechanism
- Rate limiting
- API key management

## Endpoints

### Users
```http
POST /users/register
POST /users/login
GET /users/me
```

### Transactions
```http
POST /transactions
GET /transactions
GET /transactions/{id}
```

### Accounts
```http
POST /accounts/create
GET /accounts
GET /accounts/{id}
```