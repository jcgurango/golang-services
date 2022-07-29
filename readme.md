# Ledger Application
A couple of services for double-entry bookkeeping. Code generation provide from [GrantZheng/kit](https://github.com/GrantZheng/kit)

## Ledger Service
Requires the environment variables "DB_DRIVER" (for now only "postgres" is supported), "DB_CONNECTION_STRING" which is a postgres connection string, "JWT_SECRET" which is used to encrypt JWT tokens used for authentication, and "LOGGING_SERVICE_ADDRESS" which is the address for the logging service.

### Running
This command will run the ledger service:

```bash
env DB_DRIVER=postgres "DB_CONNECTION_STRING=postgresql://localhost/" JWT_SECRET=test123 go run ledger/cmd
```

This command will run migrations to set up the DB:

```bash
env DB_DRIVER=postgres "DB_CONNECTION_STRING=postgresql://localhost/" JWT_SECRET=test123 go run ledger/cmd/migrate/main.go
```

This command will rollback migrations to set up the DB:

```bash
env DB_DRIVER=postgres "DB_CONNECTION_STRING=postgresql://localhost/" JWT_SECRET=test123 go run ledger/cmd/migrate/main.go rollback
```

### Running Tests
```bash
go test -v ledger/...
```

### Requests
By default, the HTTP server runs on *:8081. Requests return with an "error" key which has a value if there's an error. For example:

```json
{
  "error": "That user already exists in the database"
}
```

#### POST /register
Registers a new user account.

```
POST /register
Content-Type: application/json

{
  "user": "newusername",
  "pass": "password"
}
```

Response:

```json
{
  "err": null
}
```

#### POST /authenticate
Get an authentication token to be used as an `Authorization` header. These tokens are generated via JWTs.

```
POST /authenticate
Content-Type: application/json

{
  "user": "jcgurango",
  "pass": "test123"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjIwMjItMDctMjlUMTk6NTk6NDQuMTA1MTE3MTUyKzA4OjAwIiwidXNlcl9pZCI6IjEifQ.-sBuuBT5Ob01PcTOmZwvGOFabYODkF0RezZZbzijTX0",
  "err": null
}
```

#### POST /new-account
Create a new account within the user's ledger.

```
POST /authenticate
Content-Type: application/json

{
  "name": "Cash"
}
```

Response:

```json
{
  "err": null
}
```

#### POST /new-transaction
Create a new transaction within the user's ledger.

```
POST /authenticate
Content-Type: application/json

{
  "credit_account": "1",
  "debit_account": "2",
  "amount": "10000",
  "detail": "Withdrawal of cash"
}
```

Response:

```json
{
  "err": null
}
```

#### GET /get-accounts
Lists all accounts within the user's ledger.

Response:

```json
{
  "accounts": [
    {
      "id": 1,
      "name": "Cash",
      "User": 1
    },
    {
      "id": 2,
      "name": "Notes Payable",
      "User": 1
    },
    {
      "id": 3,
      "name": "Notes Receivable",
      "User": 1
    }
  ],
  "err": null
}
```

#### GET /get-balance
Lists the running balance for all accounts within the user's ledger.

Response:

```json
{
  "balances": [
    {
      "account_id": 1,
      "account_name": "Cash",
      "balance": "₱40,000.00"
    },
    {
      "account_id": 2,
      "account_name": "Notes Payable",
      "balance": "(₱40,000.00)"
    },
    {
      "account_id": 3,
      "account_name": "Notes Receivable",
      "balance": "₱0.00"
    }
  ],
  "err": null
}
```

#### GET /get-transactions
Lists the running balance for all accounts within the user's ledger.

Response:

```json
{
  "transactions": [
    {
      "ID": 1,
      "Detail": "Withdrawal of cash",
      "Amount": "₱10,000.00",
      "DebitAccount": 2,
      "CreditAccount": 1
    },
    {
      "ID": 2,
      "Detail": "Withdrawal of cash",
      "Amount": "₱10,000.00",
      "DebitAccount": 2,
      "CreditAccount": 1
    },
    {
      "ID": 3,
      "Detail": "Withdrawal of cash",
      "Amount": "₱10,000.00",
      "DebitAccount": 2,
      "CreditAccount": 1
    },
    {
      "ID": 4,
      "Detail": "Withdrawal of cash",
      "Amount": "₱10,000.00",
      "DebitAccount": 2,
      "CreditAccount": 1
    }
  ],
  "err": null
}
```
## Logging Service
The logging service doesn't require any environment variables, and accepts logs over gRPC from other services and logs them to stdout.

### Running
```bash
go run logging/cmd
```
