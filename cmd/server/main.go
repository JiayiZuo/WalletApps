// cmd/server/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"

	"WalletApps/config"
	"WalletApps/internal/common"
	"WalletApps/internal/handler"
	"WalletApps/internal/middleware"
	"WalletApps/internal/repository"
	"WalletApps/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	common.InitLogger()
	defer common.Logger.Sync()

	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DB_DSN)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.REDIS_ADDR,
	})

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	// DI
	repo := repository.NewWalletRepository(db)
	svc := service.NewWalletService(repo, redisClient)
	h := handler.NewWalletHandler(svc)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// Routes
	r.HandleFunc("/register", userHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", userHandler.LoginHandler).Methods("POST")

	walletRouter := r.PathPrefix("/api/wallet").Subrouter()
	walletRouter.Use(middleware.JWTAuthMiddleware)
	walletRouter.HandleFunc("/create", h.CreateWallet).Methods("POST")
	walletRouter.HandleFunc("/query", h.GetWallets).Methods("GET")
	walletRouter.HandleFunc("/deposit", h.Deposit).Methods("POST")
	walletRouter.HandleFunc("/withdraw", h.Withdraw).Methods("POST")
	walletRouter.HandleFunc("/transfer", h.Transfer).Methods("POST")
	walletRouter.HandleFunc("/balance/{wallet_id}", h.GetBalance).Methods("GET")
	walletRouter.HandleFunc("/transactions/{wallet_id}", h.GetTransactions).Methods("GET")

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
