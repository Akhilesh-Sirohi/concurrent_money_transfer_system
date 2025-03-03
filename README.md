# Concurrent Money Transfer System

A high-performance, concurrent money transfer system built with Go. This system allows users to create accounts, manage wallets, and perform secure money transfers between users with proper concurrency control.

## Features

- User account management
- Wallet creation and management
- Secure money transfers between users
- Concurrent transaction processing with proper locking mechanisms
- In-memory data storage with thread-safe operations
- RESTful API for all operations

## Project Structure 
├── internals/
│ ├── server/
│ │ └── router.go
│ ├── users/
│ │ ├── controller.go
│ │ ├── model.go
│ │ ├── repo.go
│ │ └── service.go
│ ├── wallet/
│ │ ├── controller.go
│ │ ├── model.go
│ │ ├── repo.go
│ │ └── service.go
│ └── transactions/
│ ├── controller.go
│ ├── model.go
│ ├── repo.go
│ └── service.go
├── utils/
│ ├── error_code.go
│ ├── response.go
│ └── validator.go
├── tests/
│ ├── transaction/
│ │ ├── transaction_test.go
│ │ └── test_data.json
│ └── test.go
└── main.go
```

## Setup and Installation

### Prerequisites

- Go 1.18 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/concurrent_money_transfer_system.git
   cd concurrent_money_transfer_system
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Application

### Start the server

```bash
go run main.go
```

### Start the server with test users

```bash
go run main.go with_test_users
```

This will create several test users with pre-loaded wallets for testing purposes.

## API Documentation

### User Management

#### Create a new user

```bash
curl --location 'http://127.0.0.1:8080/api/user/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "4",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "balance": 100.25,
    "password": "Password123!"
}'
```

#### Get all users

```bash
curl --location 'http://127.0.0.1:8080/api/user'
```

#### Get a specific user

```bash
curl --location 'http://127.0.0.1:8080/api/user/{user_id}'
```

### Transaction Management

#### Create a money transfer

```bash
curl --location 'http://127.0.0.1:8080/api/transaction/transfer' \
--header 'Content-Type: application/json' \
--data '{
    "sender_id": "1",
    "receiver_id": "2",
    "amount": 2,
    "currency": "USD"
}'
```

#### Get all transactions

```bash
curl --location 'http://127.0.0.1:8080/api/transaction'
```

#### Get a specific transaction

```bash
curl --location 'http://127.0.0.1:8080/api/transaction/{transaction_id}'
```

#### Get all transactions for a user

```bash
curl --location 'http://127.0.0.1:8080/api/transaction/user/{user_id}'
```

## Testing

Run all tests:

```bash
go test ./...
```

Run transaction tests:

```bash
go test ./tests/transaction
```

## Concurrency Model

The system uses a mutex-based locking mechanism to ensure safe concurrent access to wallet balances. When a transaction is initiated:

1. Locks are acquired on both sender and receiver wallets in a consistent order to prevent deadlocks
2. Balances are checked and updated atomically
3. Locks are released after the transaction is completed

This approach ensures that:
- No race conditions occur during balance updates
- Insufficient balance conditions are properly detected
- The system maintains consistency even under high concurrency

## Design Decisions

- **In-memory storage**: For simplicity, the system uses in-memory storage with thread-safe data structures
- **Singleton pattern**: Service and repository instances are implemented as singletons
- **Layered architecture**: The system follows a clean separation of concerns with controllers, services, and repositories
- **Explicit locking**: The system uses explicit locking rather than relying on database transactions

## Future Improvements

- Add persistent storage with a database
- Implement authentication and authorization
- Add transaction rollback mechanisms
- Implement rate limiting
- Add more comprehensive logging and monitoring