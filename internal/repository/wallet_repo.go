package repository

import (
	"database/sql"
	"errors"

	"WalletApps/internal/model"

	"github.com/google/uuid"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetWalletByUserID(userID uuid.UUID) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.QueryRow(`
		SELECT id, user_id, balance, updated_at
		FROM wallets
		WHERE user_id = $1
	`, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) UpdateWalletBalance(walletID uuid.UUID, amount float64) error {
	_, err := r.db.Exec(`
		UPDATE wallets
		SET balance = $1, updated_at = NOW()
		WHERE id = $2
	`, amount, walletID)
	return err
}

func (r *WalletRepository) InsertTransaction(tx model.Transaction) error {
	_, err := r.db.Exec(`
		INSERT INTO transactions(id, wallet_id, amount, type, description, related_user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, tx.ID, tx.WalletID, tx.Amount, tx.Type, tx.Description, tx.RelatedUserID, tx.CreatedAt)
	return err
}

func (r *WalletRepository) GetTransactions(walletID uuid.UUID) ([]model.Transaction, error) {
	rows, err := r.db.Query(`
		SELECT id, wallet_id, amount, type, description, related_user_id, created_at
		FROM transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
	`, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var relatedUserID sql.NullString
		err := rows.Scan(&t.ID, &t.WalletID, &t.Amount, &t.Type, &t.Description, &relatedUserID, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		if relatedUserID.Valid {
			uid, _ := uuid.Parse(relatedUserID.String)
			t.RelatedUserID = &uid
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *WalletRepository) Transfer(fromWalletID, toWalletID uuid.UUID, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// lock both wallets
	var fromBalance float64
	err = tx.QueryRow(`SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`, fromWalletID).Scan(&fromBalance)
	if err != nil {
		return err
	}
	if fromBalance < amount {
		return errors.New("insufficient funds")
	}

	var toBalance float64
	err = tx.QueryRow(`SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`, toWalletID).Scan(&toBalance)
	if err != nil {
		return err
	}

	// update balances
	_, err = tx.Exec(`UPDATE wallets SET balance = $1 WHERE id = $2`, fromBalance-amount, fromWalletID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE wallets SET balance = $1 WHERE id = $2`, toBalance+amount, toWalletID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
