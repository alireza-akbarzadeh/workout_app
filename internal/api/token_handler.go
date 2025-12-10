package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/alireza-akbarzadeh/fem_project/internal/store"
	"github.com/alireza-akbarzadeh/fem_project/internal/tokens"
	"github.com/alireza-akbarzadeh/fem_project/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (th *TokenHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("failed to decode create token request:%v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	user, err := th.userStore.GetUserByUserName(req.UserName)
	if err != nil {
		th.logger.Printf("failed to get user by username:%v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to fetch user"})
		return
	}
	if user == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid username or password"})
		return
	}
	match, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		th.logger.Printf("failed to check password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !match {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}
	token, err := th.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		th.logger.Printf("failed to create new token:%v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create token"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
