// cmd/server/main.go
package main

import (
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"

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

	// Routes
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	walletRouter := r.PathPrefix("/api/wallet").Subrouter()
	walletRouter.Use(middleware.JWTAuthMiddleware)
	walletRouter.HandleFunc("/deposit", h.Deposit).Methods("POST")
	walletRouter.HandleFunc("/withdraw", h.Withdraw).Methods("POST")
	walletRouter.HandleFunc("/transfer/{to_user_id}", h.Transfer).Methods("POST")
	walletRouter.HandleFunc("/balance", h.GetBalance).Methods("GET")
	walletRouter.HandleFunc("/transactions", h.GetTransactions).Methods("GET")

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
