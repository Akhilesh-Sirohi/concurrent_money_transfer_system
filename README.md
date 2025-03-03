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
```
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

### Start the server without test users at port (8080)

```bash
go run main.go
```

### Start the server with test users at port (8080)

```bash
go run main.go with_test_users
```

This will create 3 test users with pre-loaded wallets for testing purposes with details as follows:

- User 1: id = 1, First Name = Mark, Email = mark@facebook.com, Phone Number = +1234567890, balance = 100$
- User 2: id = 2, First Name = Jane, Email = jane@gmail.com, Phone Number = +1234567890, balance = 50$
- User 3: id = 3, First Name = Adam, Email = adam@gmail.com, Phone Number = +1234567890, balance = 0$



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

Run user tests:

```bash
go test ./tests/user
```

Run transaction tests:

```bash
go test ./tests/transaction
```

## Locking Strategy

The system uses a mutex-based locking mechanism to ensure safe concurrent access to wallet balances. 

## Key Components

1. **Per-Wallet Mutex**
   - Each wallet has its own dedicated mutex
   - This ensures that only one transaction can modify a wallet's balance at a time
   - Enables concurrent operations on different wallets without interfering with each other

2. **Deadlock Prevention**
   - Always acquires locks in a deterministic order (by wallet ID)
   - If sender ID < receiver ID: locks sender first, then receiver
   - If receiver ID < sender ID: locks receiver first, then sender
   - This ensures that no two transactions can wait indefinitely for each other to release their locks, preventing deadlocks

2. **Locking Management**
   - When a transaction is initiated, it acquires locks on both the sender and receiver wallets using `GetWalletForUpdate` that acquires the lock via walletMutex.(*sync.Mutex).Lock()
   - This ensures that no other transactions can modify the balances of the sender or receiver wallets until the current transaction is complete
   - Deferred calls to ReleaseGetWalletForUpdateLock ensure locks are released after the transaction is completed or if panics occur

4. **Transaction Atomicity**
   - Locks are held throughout the entire transaction process ensuring that the transaction is atomic and consistent
   - If any part of the transaction fails, the locks are released and the transaction is rolled back

This approach ensures that:
- No race conditions occur during balance updates
- No deadlocks occur
- No Double Spending occurs
- High throughput for independent transactions
- Resilience under high concurrency (tested with 50,000 concurrent transfers)

## Design Decisions

- **In-memory storage**: For simplicity, the system uses in-memory storage with thread-safe data structures
- **Singleton pattern**: Service and repository instances are implemented as singletons
- **Layered architecture**: The system follows a clean separation of concerns with controllers, services, and repositories
- **Explicit locking**: The system uses explicit locking on wallets rather than relying on database transactions

## Functional Tests Written
### User
- TestCreateUser
- TestGetUser
- TestInvalidEmailFormat
- TestInvalidPhoneFormat
- TestPasswordTooShort
- TestEmailIsRequired

### Transactions
- TestTransferMoney
- TestGetTransaction
- TestValidateNegativeAmountTransferMoneyRequest
- TestValidateSenderAndReceiverSameTransferMoneyRequest
- TestValidateSenderDoesNotExistTransferMoneyRequest
- TestCannotTransferMoreThanBalance
- TestConcurrentTransferMoney (Make 50k concurrent transfers) and validate that all transfers are sucessful, no money is lost, final wallet balance correct.

## Future Improvements

- Add persistent storage with a database
- Implement authentication and authorization
- Add transaction rollback mechanisms
- Implement rate limiting
- Add more comprehensive logging and monitoring
