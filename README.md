# Wallet Service

A RESTful wallet service built with Go, Fiber, and PostgreSQL. This service provides APIs for managing wallet transactions, including creating transactions, checking balances, and retrieving transaction history.

## Features

- Create and manage wallet transactions (credit/debit)
- Check wallet balances
- View transaction history
- Swagger API documentation
- PostgreSQL database with connection pooling
- Environment-based configuration
- Database migrations support

## Prerequisites

- Go 1.16 or higher
- PostgreSQL
- Make (for running migrations)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/arezvani/wallet-go.git
cd wallet-go
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file like `.env.example` in the root directory.


## Running the Application

1. Start the server:
```bash
make start
```

Or run directly:
```bash
go run .
```

The server will start on `http://0.0.0.0:5000`

## API Endpoints

### Transactions

#### Create Transaction
- **POST** `/api/transaction`
- Creates a new transaction for a wallet
- Request body:
```json
{
    "wallet_id": "string",
    "amount": number,
    "type": "credit" | "debit"
}
```

#### Get Transactions
- **GET** `/api/transactions/{walletId}`
- Retrieves all transactions for a specific wallet

#### Get Balance
- **GET** `/api/balance/{walletId}`
- Retrieves the current balance for a specific wallet

## API Documentation

Swagger documentation is available at:
```
http://localhost:5000/swagger
```

## Database Migrations

### Creating a New Migration
The first step is create an empty migration with the next command.

```bash
make migration-create
```

Then enter the migration name when prompted.

This command will be generated two empty files on database/migrations/ folder: 000001_$$(date +%Y%m%d%H%M%S)_$$name.up.sql and 000001_$$(date +%Y%m%d%H%M%S)_$$name.down.sql. in up file you should define sql for create tables, index and ... and in down file define sql for revert all changes made by up like drop.

### Running Migrations
```bash
# Apply migrations
make migration-up

# Rollback migrations
make migration-down
```

## Project Structure

```
.
├── app/
│   ├── controllers/     # Request handlers
│   └── models/         # Data models
├── platform/
│   └── database/       # Database configuration and migrations
├── pkg/
│   ├── configs/        # Application configuration
│   ├── middleware/     # HTTP middleware
│   └── routes/         # Route definitions
├── .env               # Environment variables
├── main.go           # Application entry point
└── Makefile          # Build and migration commands
```

## Development

### Updating Dependencies
```bash
make update-deps
```

### Building the Application
```bash
make build
```


## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
