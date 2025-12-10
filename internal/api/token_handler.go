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

// @Summary      Register a new user account
// @Description  Creates a new user in the system after validating the input.
// @Description  The user must provide a unique username, a valid email address, and a secure password.
// @Description  Optionally, a bio can be included for the user's profile.
// @Tags         Users
// @Accept       json
// @Produce      json
//
// @Param user body createTokenRequest true "User Params"
//
// @Success      201 {object} utils.Envelope{tokens=tokens.Token} "Returns created user information"
// @Failure      400 {object} utils.Envelope "Invalid input or bad request"
// @Failure      500 {object} utils.Envelope "Server error while creating the user"
//
// @Router       /tokens [post]
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
