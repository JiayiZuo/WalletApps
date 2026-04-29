package handler

import (
	"WalletApps/internal/common"
	"encoding/json"
	"net/http"

	"WalletApps/internal/middleware"
	"WalletApps/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type WalletHandler struct {
	svc *service.WalletService
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func writeJSON(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	uidStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, common.TransactionNoPermission, common.CodeNoPermission)
		return
	}
	uid, _ := uuid.Parse(uidStr)

	wallet, err := h.svc.CreateWallet(uid)
	if err != nil {
		http.Error(w, "failed to create wallet", http.StatusInternalServerError)
		return
	}
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"code": 0,
		"data": wallet,
		"msg":  "wallet created successfully",
	})
}

func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	uidStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, common.TransactionNoPermission, common.CodeNoPermission)
		return
	}
	uid, _ := uuid.Parse(uidStr)
	wallets, err := h.svc.GetWalletsByUserID(uid)
	if err != nil {
		http.Error(w, "failed to get user's wallets", http.StatusInternalServerError)
		return
	}
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"code": 0,
		"data": wallets,
		"msg":  "get wallets successfully",
	})
}

func (h *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	// Get userID from params
	//userID := mux.Vars(r)["user_id"]
	//uid, _ := uuid.Parse(userID)

	// Using jwt
	uidStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, common.TransactionNoPermission, common.CodeNoPermission)
		return
	}
	uid, _ := uuid.Parse(uidStr)

	var req struct {
		Amount   float64   `json:"amount"`
		WalletID uuid.UUID `json:"wallet_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := h.svc.Deposit(uid, req.WalletID, req.Amount)
	if err != nil {
		writeJSON(w, common.CodeInternalError, err.Error(), nil)
		return
	}
	balance, _ := h.svc.GetBalance(uid, req.WalletID)
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"user_id":   uid,
		"wallet_id": req.WalletID,
		"amount":    req.Amount,
		"balance":   balance,
	})
}

func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	//userID := mux.Vars(r)["user_id"]
	//uid, _ := uuid.Parse(userID)
	// Using jwt
	uidStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, common.TransactionNoPermission, common.CodeNoPermission)
		return
	}
	uid, _ := uuid.Parse(uidStr)
	var req struct {
		WalletID uuid.UUID `json:"wallet_id"`
		Amount   float64   `json:"amount"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := h.svc.Withdraw(uid, req.WalletID, req.Amount)
	if err != nil {
		writeJSON(w, common.CodeInsufficientFunds, err.Error(), nil)
		return
	}
	balance, _ := h.svc.GetBalance(uid, req.WalletID)
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"user_id":   uid,
		"wallet_id": req.WalletID,
		"amount":    req.Amount,
		"balance":   balance,
	})
}

func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	//fromID := mux.Vars(r)["from_user_id"]
	//fromUUID, _ := uuid.Parse(fromID)

	// Using jwt
	fromUserIDStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, common.TransactionNoPermission, common.CodeNoPermission)
		return
	}
	fromUUID, _ := uuid.Parse(fromUserIDStr)

	var req struct {
		ToUserID     uuid.UUID `json:"to_user_id"`
		FromWalletID uuid.UUID `json:"from_wallet_id"`
		ToWalletID   uuid.UUID `json:"to_wallet_id"`
		Amount       float64   `json:"amount"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := h.svc.Transfer(fromUUID, req.ToUserID, req.FromWalletID, req.ToWalletID, req.Amount)
	if err != nil {
		writeJSON(w, common.CodeInternalError, err.Error(), nil)
		return
	}
	fromBalance, _ := h.svc.GetBalance(fromUUID, req.FromWalletID)
	toBalance, _ := h.svc.GetBalance(req.ToUserID, req.ToWalletID)
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"from_user_id":   fromUUID,
		"to_user_id":     req.ToUserID,
		"amount":         req.Amount,
		"from_wallet_id": req.FromWalletID,
		"to_wallet_id":   req.ToWalletID,
		"from_balance":   fromBalance,
		"to_balance":     toBalance,
	})
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	//userID := mux.Vars(r)["user_id"]
	//uid, _ := uuid.Parse(userID)
	uidStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, common.TransactionNoPermission, common.CodeNoPermission)
		return
	}
	uid, _ := uuid.Parse(uidStr)
	walletID := mux.Vars(r)["wallet_id"]
	walletIDStr, _ := uuid.Parse(walletID)
	balance, err := h.svc.GetBalance(uid, walletIDStr)
	if err != nil {
		writeJSON(w, common.CodeInternalError, err.Error(), nil)
		return
	}
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"user_id":   uid,
		"wallet_id": walletIDStr,
		"balance":   balance,
	})
}

func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	//userID := mux.Vars(r)["user_id"]
	//uid, _ := uuid.Parse(userID)
	uidStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, common.TransactionNoPermission, common.CodeNoPermission)
		return
	}
	uid, _ := uuid.Parse(uidStr)
	walletID := mux.Vars(r)["wallet_id"]
	walletIDStr, _ := uuid.Parse(walletID)
	txs, err := h.svc.GetTransactions(uid, walletIDStr)
	if err != nil {
		writeJSON(w, common.CodeInternalError, err.Error(), nil)
		return
	}
	writeJSON(w, common.CodeOK, common.SUCCESS, txs)
}
