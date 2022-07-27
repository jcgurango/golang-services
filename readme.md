# Ledger Application
A couple of services for double-entry bookkeeping. Code generation provide from [GrantZheng/kit](https://github.com/GrantZheng/kit)

## Ledger Service
Requires the environment variables "DB_DRIVER" (for now only "postgres" is supported), "DB_CONNECTION_STRING" which is a postgres connection string, and "JWT_SECRET" which is used to encrypt JWT tokens used for authentication.

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
