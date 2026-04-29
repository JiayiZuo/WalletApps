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
│ │ ├── utils.go
│ │ ├── jwt.go
│ │ └── logger.go
│ ├── handler
│ │ ├── user_handler.go
│ │ └── wallet_handler.go
│ ├── middleware
│ │ ├── jwt.go
│ │ └── logger.go
│ ├── model
│ │ └── models.go
│ ├── repository
│ │ ├── user_repo.go
│ │ └── wallet_repo.go
│ ├── service
│ │ ├── user_service.go
│ │ └── wallet_service.go
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
| Method | Endpoint                              | Description           |
|--------|---------------------------------------|-----------------------|
| POST   | `/register`                           | Register new account  |
| POST   | `/login`                              | Login & get JWT token |
| POST   | `/api/wallet/create`                  | Create new wallet     |
| GET    | `/api/wallet/query  `                 | Query user's wallets  |
| POST   | `/api/wallet/deposit`                 | Deposit funds         |
| POST   | `/api/wallet/withdraw`                | Withdraw funds        |
| POST   | `/api/wallet/transfer `               | Transfer funds        |
| GET    | `/api/wallet/balance/{wallet_id}`     | Get balance           |
| GET    | `/api/wallet/transactions/{wallet_id}`| Get transactions      |

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
    "username": "Jiayi",
    "password": "123456"
}
// response
{
    "name": "Jiayi",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzc1MzE1NzUsInVzZXJfaWQiOiI3NGRjZDRmOS1jOTQ0LTQ5MmQtODAyOS1jNjVjMDA4MjQ5Y2UifQ.puHbIbc7TQ0UAx1drLQg5bzGgLHMKq-lik3zpS5Y3yk",
    "user_id": "74dcd4f9-c944-492d-8029-c65c008249ce"
}
````
2. Register
````
request:
// request http://127.0.0.1:8080/register POST
{
    "username": "Edward",
    "password": "123456"
}
//response
{
    "code": 0,
    "msg": "Registration successful",
    "data": {
        "user_id": "9a46377c-bad1-4eb2-9511-511f0cda7f90",
        "username": "Edward"
    }
}
````
3. Deposit
````
// request http://127.0.0.1:8080/api/wallet/deposit POST Authorization: Bearer Token ••••••
{
    "wallet_id": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8",
    "amount": 50.50
}
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "amount": 50.5,
        "balance": 388.75,
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
		"wallet_id": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8"
    }
}
````
4. Withdraw
````
// request http://127.0.0.1:8080/api/wallet/withdraw POST Authorization: Bearer Token ••••••
{
    "wallet_id": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8",
    "amount": 10.25
}
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "amount": 10.25,
        "balance": 349.25,
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
		"wallet_id": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8"
    }
}
````
5. Transfer
````
// request http://127.0.0.1:8080/api/wallet/transfer POST Authorization: Bearer Token ••••••
{
    "to_user_id": "9a46377c-bad1-4eb2-9511-511f0cda7f90",
    "from_wallet_id": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8",
    "to_wallet_id": "107ebac4-d9a9-4a55-b56f-b930e0d7d110",
    "amount": 100
}
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "amount": 100,
        "from_balance": 700,
        "from_user_id": "74dcd4f9-c944-492d-8029-c65c008249ce",
        "from_wallet_id": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8",
        "to_balance": 100,
        "to_user_id": "9a46377c-bad1-4eb2-9511-511f0cda7f90",
        "to_wallet_id": "107ebac4-d9a9-4a55-b56f-b930e0d7d110"
    }
}
````
6. Balance
````
// request http://127.0.0.1:8080/api/wallet/balance/e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8 GET Authorization: Bearer Token ••••••
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "balance": 338.25,
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
		"wallet_id": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8"
    }
}
````
7. Transactions
````
// request http://127.0.0.1:8080/api/wallet/transactions/e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8 Authorization: Bearer Token ••••••
// response
{
    "code": 0,
    "msg": "success",
    "data": [
        {
            "ID": "ee42275b-5b8b-459b-9a4b-b448fb427d34",
            "WalletID": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8",
            "Amount": -200,
            "Type": "withdraw",
            "Description": "",
            "RelatedUserID": null,
            "CreatedAt": "2026-04-29T06:33:48.181228Z"
        },
        {
            "ID": "2ed7d9d8-7e75-4e2b-bbe0-d34316eeec29",
            "WalletID": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8",
            "Amount": 1000,
            "Type": "deposit",
            "Description": "",
            "RelatedUserID": null,
            "CreatedAt": "2026-04-29T04:00:22.685992Z"
        }
    ]
}
````
8. CreateWallet
````
// request http://127.0.0.1:8080/api/wallet/create POST Authorization: Bearer Token ••••••
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "code": 0,
        "data": {
            "ID": "107ebac4-d9a9-4a55-b56f-b930e0d7d110",
            "UserID": "9a46377c-bad1-4eb2-9511-511f0cda7f90",
            "Address": "0x13237c478e9dd1eca7da422e373991af7cdce2a3",
            "Balance": 0,
            "UpdatedAt": "2026-04-29T14:41:54.5479192+08:00"
        },
        "msg": "wallet created successfully"
    }
}
````
9. GetWallets
````
// request http://127.0.0.1:8080/api/wallet/query GET Authorization: Bearer Token ••••••
// response
{
    "code": 0,
    "msg": "success",
    "data": {
        "code": 0,
        "data": [
            {
                "ID": "e575aeeb-91aa-4089-9b8d-1e4daa6dcdc8",
                "UserID": "74dcd4f9-c944-492d-8029-c65c008249ce",
                "Address": "0xfc2e6075bf081952735797e36ad072f417ce9993",
                "Balance": 700,
                "UpdatedAt": "2026-04-29T06:33:48.175034Z"
            }
        ],
        "msg": "get wallets successfully"
    }
}
````