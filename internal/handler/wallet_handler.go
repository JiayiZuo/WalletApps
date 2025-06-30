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

// LoginHandler A simple login api to show jwt demo
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// mock user
	var req struct {
		UserID string `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.UserID == "" {
		http.Error(w, "missing user_id", common.CodeInvalidParam)
		return
	}
	token, _ := common.GenerateJWT(req.UserID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
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
		Amount float64 `json:"amount"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := h.svc.Deposit(uid, req.Amount)
	if err != nil {
		writeJSON(w, common.CodeInternalError, err.Error(), nil)
		return
	}
	balance, _ := h.svc.GetBalance(uid)
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"user_id": uid,
		"amount":  req.Amount,
		"balance": balance,
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
		Amount float64 `json:"amount"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := h.svc.Withdraw(uid, req.Amount)
	if err != nil {
		writeJSON(w, common.CodeInsufficientFunds, err.Error(), nil)
		return
	}
	balance, _ := h.svc.GetBalance(uid)
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"user_id": uid,
		"amount":  req.Amount,
		"balance": balance,
	})
}

func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	//fromID := mux.Vars(r)["from_user_id"]
	//fromUUID, _ := uuid.Parse(fromID)

	// Using jwt
	fromIDStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		http.Error(w, common.TransactionNoPermission, common.CodeNoPermission)
		return
	}
	fromUUID, _ := uuid.Parse(fromIDStr)
	toID := mux.Vars(r)["to_user_id"]
	toUUID, _ := uuid.Parse(toID)

	var req struct {
		Amount float64 `json:"amount"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := h.svc.Transfer(fromUUID, toUUID, req.Amount)
	if err != nil {
		writeJSON(w, common.CodeInternalError, err.Error(), nil)
		return
	}
	fromBalance, _ := h.svc.GetBalance(fromUUID)
	toBalance, _ := h.svc.GetBalance(toUUID)
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"from_user_id": fromUUID,
		"to_user_id":   toUUID,
		"amount":       req.Amount,
		"from_balance": fromBalance,
		"to_balance":   toBalance,
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
	balance, err := h.svc.GetBalance(uid)
	if err != nil {
		writeJSON(w, common.CodeInternalError, err.Error(), nil)
		return
	}
	writeJSON(w, common.CodeOK, common.SUCCESS, map[string]interface{}{
		"user_id": uid,
		"balance": balance,
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
	txs, err := h.svc.GetTransactions(uid)
	if err != nil {
		writeJSON(w, common.CodeInternalError, err.Error(), nil)
		return
	}
	writeJSON(w, common.CodeOK, common.SUCCESS, txs)
}
