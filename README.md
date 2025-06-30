# A simple wallet management service built with Golang:
````
✅ Featuring JWT-based authentication
✅ PostgreSQL storage
✅ Redis lock to control concurrency
✅ Global log config including api and service(Can be used to monitor traffic and business)
✅ Unit test to validate logical correctness
````

# Implemented Features
````
✅ Management of user wallets(deposit, withdraw, transfer, get balance, get transactions)
✅ Secure JWT-based authentication
✅ RESTful API
✅ Modular and testable Go code
✅ Database migrations
````

# Project Structure
````
├── README.md
├── cmd
│ └── server
│ └── main.go
├── config
│ └── config.go
├── go.mod
├── go.sum
├── internal
│ ├── common
│ │ ├── code.go
│ │ ├── constant.go
│ │ ├── jwt.go
│ │ └── logger.go
│ ├── handler
│ │ └── wallet_handler.go
│ ├── middleware
│ │ ├── jwt.go
│ │ └── logger.go
│ ├── model
│ │ └── models.go
│ ├── repository
│ │ └── wallet_repo.go
│ └── service
│ └── wallet_service.go
├── migrations
│ └── 001_init.sql
├── tests
  └── wallet_test.go
````

# Getting Start
Prerequisites
````
Go 1.22+
PostgreSQL
Redis
````
# Installation
1. Clone the repository:
````
git clone https://github.com/zuojiayi/WalletApps.git
cd wallet-service
````
2. Install dependencies:
````
go mod tidy
````
3. Setup your .env file (or update config variables in config.go) with your database connection info:env
````
DB_DSN=postgres://user:password@localhost:5432/wallet?sslmode=disable
JWT_SECRET=your_jwt_secret
````
4. Run migrations:
````
psql -U youruser -d wallet -f migrations/001_init.sql
````
5. Start the server:
````
go run cmd/server/main.go
````

# API
| Method | Endpoint                            | Description           |
|--------|-------------------------------------|-----------------------|
| POST   | `/login`                            | Login & get JWT token |
| POST   | `/api/wallet/deposit`               | Deposit funds         |
| POST   | `/api/wallet/withdraw`              | Withdraw funds        |
| POST   | `/api/wallet/transfer/{to_user_id}` | Transfer funds        |
| GET    | `/api/wallet/balance`               | Get balance           |
| GET    | `/api/wallet/transactions`          | Get transactions      |

# Testing
````
go test ./tests/...
````

# How to view the code
````
1. API routing: defined in cmd/server/main.go
2. Request handlers: in internal/handler/wallet_handler.go
3. Business logic: in internal/service/wallet_service.go
4. PostgreSQL schema: in migrations/001_init.sql
5. Database models: in internal/model/models.go; SQL queries in internal/repository/wallet_repo.go
6. JWT auth: all requests require JWT; use /login to get a token and add it to the Authorization header as a Bearer token
7. JWT and logging: implemented in internal/common and internal/middleware
8. Global response messages: in internal/common/constant.go, business codes in internal/common/code.go
````

# Areas to be improved
````
1. Add Swagger or OpenAPI documentation for easier client integration
2. Add Dockerfile and docker-compose for one-click deployment
3. Add rate limiting to protect against brute-force or denial-of-service attacks
4. Implement user registration flow (currently only login is supported)
5. Support for multi-currency wallets or fiat on-ramp integration
6. Add integration tests (currently only unit tests)
7. Improve test coverage for edge cases (e.g. invalid amounts, concurrency conflicts)
8. Implement an admin panel or monitoring dashboard
9. Add CI/CD pipelines (GitHub Actions, GitLab CI, etc.)
````

# Time spent on this Project
````
7 hours for logic implement, optmization and testing by postman
1 hour for complete the README file
````

# Features not to do in the submission
all the file written in the .gitignore file:
````
.env
/bin/
/vendor/
/*.exe
*.log
*.out
*.test
.idea/
.vscode/
*.swp
.DS_Store
````

# Request and Response Demo
1. Login
request:
````
// request http://127.0.0.1:8080/login POST
{
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
// response
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTEzNjQ2ODcsInVzZXJfaWQiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAifQ.2-oEG-G5C9XIh9Ioe_J7Zb95pKeYPMIHsGhb7viEx_k"
}
````
2. Deposit
````
// request http://127.0.0.1:8080/api/wallet/deposit POST Authorization: Bearer Token ••••••
{
    "amount": 50.50
}
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "amount": 50.5,
        "balance": 388.75,
        "user_id": "550e8400-e29b-41d4-a716-446655440000"
    }
}
````
3. Withdraw
````
// request http://127.0.0.1:8080/api/wallet/withdraw POST Authorization: Bearer Token ••••••
{
    "amount": 10.25
}
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "amount": 10.25,
        "balance": 349.25,
        "user_id": "550e8400-e29b-41d4-a716-446655440000"
    }
}
````
4. Transfer
````
// request http://127.0.0.1:8080/api/wallet/transfer/{to_user_id} POST Authorization: Bearer Token ••••••
{
    "amount": 11.00
}
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "amount": 11,
        "from_balance": 338.25,
        "from_user_id": "550e8400-e29b-41d4-a716-446655440000",
        "to_balance": 102,
        "to_user_id": "dfd6d470-e316-c52c-ea3b-58b989569f5e"
    }
}
````
5. Balance
````
// request http://127.0.0.1:8080/api/wallet/balance GET Authorization: Bearer Token ••••••
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "balance": 338.25,
        "user_id": "550e8400-e29b-41d4-a716-446655440000"
    }
}
````
6. Transactions
````
// request http://127.0.0.1:8080/api/wallet/transactions Authorization: Bearer Token ••••••
// response
{
    "code": 0,
    "msg": "success",
    "data": [
        {
            "ID": "60110852-76c8-4793-95db-0f6a63552744",
            "WalletID": "660e8400-e29b-41d4-a716-446655440000",
            "Amount": 200.5,
            "Type": "deposit",
            "Description": "",
            "RelatedUserID": null,
            "CreatedAt": "2025-06-30T16:53:38.560819+08:00"
        },
        {
            "ID": "971e0632-6a2f-4702-a2db-6d2f13a28250",
            "WalletID": "660e8400-e29b-41d4-a716-446655440000",
            "Amount": -1.25,
            "Type": "withdraw",
            "Description": "",
            "RelatedUserID": null,
            "CreatedAt": "2025-06-30T16:53:25.026253+08:00"
        }
    ]
}
````