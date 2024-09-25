package handlers

import (
	"bookstore_api/internal/services"
	"bookstore_api/models"
	"bookstore_api/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type UserHandler struct {
	*Handler
	userService    *services.UserService
	sessionService *services.SessionService
}

func NewUserHandler(handler *Handler, userService *services.UserService, sessionService *services.SessionService) *UserHandler {
	return &UserHandler{
		Handler:        handler,
		userService:    userService,
		sessionService: sessionService,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userRegister := &models.UserRegister{}
	if err := json.NewDecoder(r.Body).Decode(userRegister); err != nil {
		tools.RespondWithError(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	createdUser, err := h.userService.RegisterUser(r.Context(), userRegister)
	if err != nil {
		tools.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Redirect to login page
	tools.RespondWithJSON(w, createdUser, http.StatusCreated)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	userLogin := &models.UserLogin{}
	if err := json.NewDecoder(r.Body).Decode(userLogin); err != nil {
		tools.RespondWithError(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	userResponse, err := h.userService.LoginUser(r.Context(), userLogin)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	accessClaims, accessToken, err := tools.GenerateToken(userResponse.Email, userResponse.IsAdmin, 30*time.Minute)
	if err != nil {
		tools.RespondWithError(w, errors.New("failed to generate tokens"), http.StatusInternalServerError)
		return
	}

	sessionClaims, sessionToken, err := tools.GenerateToken(userResponse.Email, userResponse.IsAdmin, 24*time.Hour)
	if err != nil {
		tools.RespondWithError(w, errors.New("failed to generate sessions"), http.StatusInternalServerError)
		return
	}

	session := &models.Sessions{
		ID:           sessionClaims.RegisteredClaims.ID,
		UserEmail:    userResponse.Email,
		RefreshToken: sessionToken,
		IsRevoked:    false,
		CreatedAt:    sessionClaims.IssuedAt.Time,
		ExpiresAt:    sessionClaims.ExpiresAt.Time,
	}

	_, err = h.sessionService.CreateSession(r.Context(), session)
	if err != nil {
		tools.RespondWithError(w, errors.New("failed to generate sessions"), http.StatusInternalServerError)
		return
	}

	loginResponse := models.UserLoginResponse{
		User:       *userResponse,
		SessionsId: sessionClaims.RegisteredClaims.ID,

		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.ExpiresAt.Time,

		RefreshToken:          sessionToken,
		RefreshTokenExpiresAt: sessionClaims.ExpiresAt.Time,
	}

	w.Header().Set("Authorization", "Bearer "+loginResponse.AccessToken)
	tools.RespondWithJSON(w, loginResponse, http.StatusOK)
}

func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	authHeader, err := getBearerToken(r)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	claims, err := tools.ValidateToken(authHeader)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	err = h.sessionService.DeleteSession(r.Context(), claims.RegisteredClaims.ID)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	authHeader, err := getBearerToken(r)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), "token", authHeader)

	userData := &models.UserUpdateData{}
	if err := json.NewDecoder(r.Body).Decode(userData); err != nil {
		tools.RespondWithError(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.userService.UpdateUserData(ctx, userData)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	tools.RespondWithJSON(w, updatedUser, http.StatusOK)
}

func (h *UserHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	authHeader, err := getBearerToken(r)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	claims, err := tools.ValidateToken(authHeader)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Use claims.RegisteredClaims.ID instead of claims.ID
	session, err := h.sessionService.GetSession(r.Context(), claims.RegisteredClaims.ID)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid session"), http.StatusUnauthorized)
		return
	}

	if claims.Subject != session.UserEmail {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	if session.IsRevoked {
		tools.RespondWithError(w, errors.New("revoked session"), http.StatusUnauthorized)
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		err = h.sessionService.RevokeSession(r.Context(), session.ID)
		if err != nil {
			tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
			return
		}
		tools.RespondWithError(w, errors.New("session is expired"), http.StatusUnauthorized)
		return
	}

	_, accessToken, err := tools.GenerateToken(session.UserEmail, claims.IsAdmin, 15*time.Minute)
	if err != nil {
		tools.RespondWithError(w, errors.New("failed to generate tokens"), http.StatusInternalServerError)
		return
	}

	accessTokenResp := &models.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: claims.ExpiresAt.Time,
	}

	w.Header().Set("Authorization", "Bearer "+accessToken)
	tools.RespondWithJSON(w, accessTokenResp, http.StatusOK)
}

func (h *UserHandler) RevokeAccessToken(w http.ResponseWriter, r *http.Request) {
	authHeader, err := getBearerToken(r)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	claims, err := tools.ValidateToken(authHeader)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	err = h.sessionService.RevokeSession(r.Context(), claims.RegisteredClaims.ID)
	if err != nil {
		tools.RespondWithError(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no Authorization header found")
	}

	// The Authorization header should start with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization header is not of type Bearer")
	}

	// Extract the token by trimming the "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}
