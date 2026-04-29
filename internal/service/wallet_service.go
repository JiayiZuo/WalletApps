package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"WalletApps/internal/model"
	"WalletApps/internal/repository"

	"WalletApps/internal/common"

	"github.com/google/uuid"
)

type WalletService struct {
	repo   *repository.WalletRepository
	locker *redislock.Client
}

func NewWalletService(r *repository.WalletRepository, redisClient *redis.Client) *WalletService {
	locker := redislock.New(redisClient)
	return &WalletService{
		repo:   r,
		locker: locker,
	}
}

func (s *WalletService) CreateWallet(userID uuid.UUID) (*model.Wallet, error) {
	address, err := common.MockGenerateAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mock address: %v", err)
	}

	wallet := &model.Wallet{
		ID:        uuid.New(),
		UserID:    userID,
		Address:   address,
		Balance:   0,
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) GetWalletsByUserID(userID uuid.UUID) ([]*model.Wallet, error) {
	wallets, err := s.repo.GetWalletsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (s *WalletService) Deposit(userID uuid.UUID, walletID uuid.UUID, amount float64) error {
	ctx := context.Background()
	lock, err := s.locker.Obtain(ctx, common.WalletLockPrefix+userID.String(), 5*time.Second, nil)
	if err == redislock.ErrNotObtained {
		return errors.New(common.CurrentDepositInProgress)
	}
	if err != nil {
		return err
	}
	defer lock.Release(ctx)

	wallet, err := s.repo.GetWalletByWalletID(walletID)
	if err != nil {
		return err
	}
	if wallet.UserID != userID {
		return errors.New("User and wallet didn't match, please check.")
	}
	newBalance := wallet.Balance + amount
	err = s.repo.UpdateWalletBalance(wallet.ID, newBalance)
	if err != nil {
		return err
	}
	t := model.Transaction{
		ID:        uuid.New(),
		WalletID:  wallet.ID,
		Amount:    amount,
		Type:      common.TransactionTypeDeposit,
		CreatedAt: time.Now(),
	}
	common.Logger.Info(common.DepositRequest,
		zap.String("user_id", userID.String()),
		zap.Float64("amount", amount),
	)
	return s.repo.InsertTransaction(t)
}

func (s *WalletService) Withdraw(userID uuid.UUID, walletID uuid.UUID, amount float64) error {
	ctx := context.Background()
	lock, err := s.locker.Obtain(ctx, common.WalletLockPrefix+userID.String(), 5*time.Second, nil)
	if err == redislock.ErrNotObtained {
		return errors.New(common.CurrentDepositInProgress)
	}
	if err != nil {
		return err
	}
	defer lock.Release(ctx)

	wallet, err := s.repo.GetWalletByWalletID(walletID)
	if err != nil {
		return err
	}
	if wallet.UserID != userID {
		return errors.New("User and wallet didn't match, please check.")
	}
	if wallet.Balance < amount {
		return errorInsufficientFunds{}
	}
	newBalance := wallet.Balance - amount
	err = s.repo.UpdateWalletBalance(wallet.ID, newBalance)
	if err != nil {
		return err
	}
	t := model.Transaction{
		ID:        uuid.New(),
		WalletID:  wallet.ID,
		Amount:    -amount,
		Type:      common.TransactionTypeWithdraw,
		CreatedAt: time.Now(),
	}
	common.Logger.Info(common.WithdrawRequest,
		zap.String("user_id", userID.String()),
		zap.Float64("amount", amount),
	)
	return s.repo.InsertTransaction(t)
}

func (s *WalletService) Transfer(fromUserID, toUserID uuid.UUID, fromWalletID, toWalletID uuid.UUID, amount float64) error {
	ctx := context.Background()

	// lock 2 wallet in case of deadlock order lock according to userID in asc order
	ids := []uuid.UUID{fromUserID, toUserID}
	if ids[0].String() > ids[1].String() {
		ids[0], ids[1] = ids[1], ids[0]
	}

	lock1, err := s.locker.Obtain(ctx, common.WalletLockPrefix+ids[0].String(), 5*time.Second, nil)
	if err == redislock.ErrNotObtained {
		return errors.New(common.CurrentDepositInProgress)
	}
	if err != nil {
		return err
	}
	defer lock1.Release(ctx)

	lock2, err := s.locker.Obtain(ctx, common.WalletLockPrefix+ids[1].String(), 5*time.Second, nil)
	if err == redislock.ErrNotObtained {
		return errors.New(common.CurrentDepositInProgress)
	}
	if err != nil {
		lock1.Release(ctx)
		return err
	}
	defer lock2.Release(ctx)

	fromWallet, err := s.repo.GetWalletByWalletID(fromWalletID)
	if err != nil {
		return err
	}
	toWallet, err := s.repo.GetWalletByWalletID(toWalletID)
	if err != nil {
		return err
	}
	if fromWallet.UserID != fromUserID || toWallet.UserID != toUserID {
		return errors.New("User and wallet didn't match, please check.")
	}
	err = s.repo.Transfer(fromWallet.ID, toWallet.ID, amount)
	if err != nil {
		return err
	}

	// record transaction from
	t1 := model.Transaction{
		ID:            uuid.New(),
		WalletID:      fromWallet.ID,
		Amount:        -amount,
		Type:          common.TransactionTypeTransfer,
		Description:   common.TransactionTypeTransfer + " to " + toUserID.String(),
		RelatedUserID: &toUserID,
		CreatedAt:     time.Now(),
	}
	_ = s.repo.InsertTransaction(t1)

	// record transaction to
	t2 := model.Transaction{
		ID:            uuid.New(),
		WalletID:      toWallet.ID,
		Amount:        amount,
		Type:          common.TransactionTypeTransfer,
		Description:   common.TransactionTypeTransfer + " from " + fromUserID.String(),
		RelatedUserID: &fromUserID,
		CreatedAt:     time.Now(),
	}
	_ = s.repo.InsertTransaction(t2)

	common.Logger.Info(common.WithdrawRequest,
		zap.String("from_id", fromUserID.String()),
		zap.String("to_id", toUserID.String()),
		zap.Float64("amount", amount),
	)
	return nil
}

func (s *WalletService) GetBalance(userID, walletID uuid.UUID) (float64, error) {
	wallet, err := s.repo.GetWalletByWalletID(walletID)
	if err != nil {
		return 0, err
	}
	if wallet.UserID != userID {
		return 0, errors.New("User and wallet didn't match, please check.")
	}
	common.Logger.Info(common.WithdrawRequest,
		zap.String("from_id", userID.String()),
		zap.Float64("amount", wallet.Balance),
	)
	return wallet.Balance, nil
}

func (s *WalletService) GetTransactions(userID, walletID uuid.UUID) ([]model.Transaction, error) {
	wallet, err := s.repo.GetWalletByWalletID(walletID)
	if err != nil {
		return nil, err
	}
	if wallet.UserID != userID {
		return nil, errors.New("User and wallet didn't match, please check.")
	}
	common.Logger.Info(common.WithdrawRequest,
		zap.String("user_id", userID.String()),
		zap.Float64("balance", wallet.Balance),
	)
	return s.repo.GetTransactions(wallet.ID)
}

type errorInsufficientFunds struct{}

func (e errorInsufficientFunds) Error() string {
	return common.TransactionTransferInsufficientFunds
}
