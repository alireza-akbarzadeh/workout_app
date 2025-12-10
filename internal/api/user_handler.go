package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/alireza-akbarzadeh/fem_project/internal/store"
	"github.com/alireza-akbarzadeh/fem_project/internal/utils"
	"github.com/alireza-akbarzadeh/fem_project/internal/validation"
)

type UserHandler struct {
	UserStore store.UserStore
	logger    *log.Logger
}

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) RegisterUser(req *registerUserRequest) error {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return errors.New("username, email, and password are required")
	}

	if len(req.Username) > 50 {
		return errors.New("username exceeds maximum length of 50 characters")
	}

	if len(req.Email) > 100 {
		return errors.New("email exceeds maximum length of 100 characters")
	}
	if !validation.IsEmailValid(req.Email) {
		return errors.New("invalid email format")
	}

	if !validation.IsPasswordValid(req.Password) {
		return errors.New("password must be at least 8 characters long and include uppercase, lowercase, number, and special character")
	}
	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("failed to decode register user request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = uh.RegisterUser(&req)
	if err != nil {
		uh.logger.Printf("validation error: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	newUser := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		newUser.Bio = req.Bio
	}

	err = newUser.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("failed to hash password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to process password"})
		return
	}

	createdUser, err := uh.UserStore.CreateUser(newUser)
	if err != nil {
		uh.logger.Printf("failed to create user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create user"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": createdUser})

}
func (uh *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username is required"})
		return
	}

	user, err := uh.UserStore.GetUserByUsername(username)
	if err != nil {
		uh.logger.Printf("failed to get user by username: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to fetch user"})
		return
	}
	if user == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})

}

func (uh *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req store.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("failed to decode update user request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = uh.UserStore.UpdateUser(&req)
	if err != nil {
		uh.logger.Printf("failed to update user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update user"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user updated successfully"})
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "user ID is required"})
		return
	}

	userID, err := utils.ParseInt64(id)
	if err != nil {
		uh.logger.Printf("invalid user ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user ID"})
		return
	}

	user, err := uh.UserStore.GetUserByID(userID)
	if err != nil {
		uh.logger.Printf("failed to get user by ID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to fetch user"})
		return
	}
	if user == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

func (uh *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "user ID is required"})
		return
	}
	userID, err := utils.ParseInt64(id)
	if err != nil {
		uh.logger.Printf("invalid user ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user ID"})
		return
	}

	err = uh.UserStore.DeleteUser(userID)
	if err != nil {
		uh.logger.Printf("failed to delete user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to delete user"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user deleted successfully"})
}
