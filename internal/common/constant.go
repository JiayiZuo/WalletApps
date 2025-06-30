package common

const (
	// golable msg
	SUCCESS = "success"

	// transaction related
	TransactionTypeDeposit               = "deposit"
	TransactionTypeWithdraw              = "withdraw"
	TransactionTypeTransfer              = "transfer"
	TransactionTransferInsufficientFunds = "insufficient funds"
	TransactionNoPermission              = "unauthorized"

	// concurrency
	WalletLockPrefix         = "lock:wallet:"
	CurrentDepositInProgress = "concurrent deposit in progress"

	// log keywords
	DepositRequest  = "deposit request"
	WithdrawRequest = "withdraw request"
)
