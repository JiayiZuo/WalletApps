package tests

import (
	"database/sql"
	"testing"

	"WalletApps/internal/repository"
	"WalletApps/internal/service"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://zuojiayi@localhost:5432/wallet?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

func TestDepositAndWithdraw(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Create test user + wallet
	userID := uuid.New()
	walletID := uuid.New()

	_, err := db.Exec(`INSERT INTO users(id,name) VALUES($1,$2)`, userID, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO wallets(id,user_id,balance) VALUES($1,$2,$3)`, walletID, userID, 0)
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewWalletRepository(db)
	svc := service.NewWalletService(repo)

	// Deposit 100
	err = svc.Deposit(userID, 100)
	if err != nil {
		t.Fatal(err)
	}

	balance, err := svc.GetBalance(userID)
	if err != nil {
		t.Fatal(err)
	}
	if balance != 100 {
		t.Fatalf("Expected balance 100 got %f", balance)
	}

	// Withdraw 40
	err = svc.Withdraw(userID, 40)
	if err != nil {
		t.Fatal(err)
	}

	balance, err = svc.GetBalance(userID)
	if err != nil {
		t.Fatal(err)
	}
	if balance != 60 {
		t.Fatalf("Expected balance 60 got %f", balance)
	}
}

func TestTransfer(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Create test users
	user1 := uuid.New()
	user2 := uuid.New()
	wallet1 := uuid.New()
	wallet2 := uuid.New()

	_, err := db.Exec(`INSERT INTO users(id,name) VALUES($1,$2)`, user1, "John")
	_, err = db.Exec(`INSERT INTO users(id,name) VALUES($1,$2)`, user2, "Lucy")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO wallets(id,user_id,balance) VALUES($1,$2,$3)`, wallet1, user1, 100)
	_, err = db.Exec(`INSERT INTO wallets(id,user_id,balance) VALUES($1,$2,$3)`, wallet2, user2, 50)
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewWalletRepository(db)
	svc := service.NewWalletService(repo)

	// Alice transfer 30 to Bob
	err = svc.Transfer(user1, user2, 30)
	if err != nil {
		t.Fatal(err)
	}

	balance1, _ := svc.GetBalance(user1)
	balance2, _ := svc.GetBalance(user2)

	if balance1 != 70 {
		t.Fatalf("Expected Alice 70 got %f", balance1)
	}
	if balance2 != 80 {
		t.Fatalf("Expected Bob 80 got %f", balance2)
	}
}
