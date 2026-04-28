package handler

import (
	"WalletApps/internal/common"
	"WalletApps/internal/service"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)
}

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var RegisterReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&RegisterReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if RegisterReq.Username == "" || RegisterReq.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	uid := uuid.New()

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		writeJSON(w, common.CodeInternalError, "Failed to generate salt", nil)
		return
	}

	saltStr := hex.EncodeToString(salt)

	raw := RegisterReq.Password + saltStr
	hash := sha256.Sum256([]byte(raw))
	passwordHash := hex.EncodeToString(hash[:])

	err := h.svc.CreateUser(uid, RegisterReq.Username, passwordHash, saltStr)
	if err != nil {
		writeJSON(w, common.CodeInternalError, "Failed to register user: "+err.Error(), nil)
		return
	}

	writeJSON(w, common.CodeOK, "Registration successful", map[string]interface{}{
		"user_id":  uid.String(),
		"username": RegisterReq.Username,
	})
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !common.CheckPassword(req.Password, user.Salt, user.PasswordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := common.GenerateJWT(user.ID.String())
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"user_id": user.ID.String(),
		"name":    user.Name,
		"token":   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
