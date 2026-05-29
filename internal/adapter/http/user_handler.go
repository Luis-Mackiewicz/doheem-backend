package http

import (
	"errors"
	"net/http"
	"time"

	"doheem-backend/internal/domain"
)

type UserHandler struct {
	svc *domain.UserService
	jwt *JWTService
}

func NewUserHandler(svc *domain.UserService, jwt *JWTService) *UserHandler {
	return &UserHandler{svc: svc, jwt: jwt}
}

// Register registers a new user
// @Summary Register a new user
// @Description Create a new user account and return an auth token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{name=string,email=string,password=string} true "Registration details"
// @Success 201 {object} authResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 409 {object} map[string]any "Email already in use"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/auth/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string  `json:"name"      validate:"required"`
		Email     string  `json:"email"     validate:"required,email"`
		Password  string  `json:"password"  validate:"required,min=6"`
		AvatarURL *string `json:"avatar_url,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	user, err := h.svc.Register(r.Context(), domain.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		AvatarURL:    req.AvatarURL,
	})
	if err != nil {
		handleError(w, err)
		return
	}
	token, err := h.jwt.GenerateToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	refreshToken, refreshHash, err := h.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}
	if err := h.svc.StoreRefreshToken(r.Context(), user.ID, refreshHash, time.Now().Add(168*time.Hour)); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store refresh token")
		return
	}
	respondJSON(w, http.StatusCreated, authResponse{
		User:         toUserResponse(user),
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// Login authenticates a user
// @Summary Login
// @Description Authenticate a user with email and password and return an auth token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string} true "Login credentials"
// @Success 200 {object} authResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Invalid email or password"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	user, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	token, err := h.jwt.GenerateToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	refreshToken, refreshHash, err := h.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}
	if err := h.svc.StoreRefreshToken(r.Context(), user.ID, refreshHash, time.Now().Add(168*time.Hour)); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store refresh token")
		return
	}
	respondJSON(w, http.StatusOK, authResponse{
		User:         toUserResponse(user),
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// GetProfile returns the authenticated user's profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} userResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "User not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/users/me [get]
// @Security BearerAuth
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	user, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toUserResponse(user))
}

// UpdateProfile updates the authenticated user's profile
// @Summary Update user profile
// @Description Update the profile of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{name=string,email=string,avatar_url=string} true "Profile update details"
// @Success 200 {object} userResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/users/me [put]
// @Security BearerAuth
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	var req struct {
		Name      *string `json:"name,omitempty"`
		Email     *string `json:"email,omitempty"     validate:"omitempty,email"`
		AvatarURL *string `json:"avatar_url,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	user, err := h.svc.Update(r.Context(), userID, domain.UpdateUserParams{
		Name:      req.Name,
		Email:     req.Email,
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		handleError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toUserResponse(user))
}

// ChangePassword changes the authenticated user's password
// @Summary Change password
// @Description Change the password of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{old_password=string,new_password=string} true "Password change details"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Invalid current password"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/users/me/password [put]
// @Security BearerAuth
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	if err := h.svc.UpdatePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "invalid current password")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Refresh returns new access and refresh tokens using a valid refresh token
// @Summary Refresh tokens
// @Description Get a new access token and refresh token pair using a valid refresh token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "Refresh token"
// @Success 200 {object} authResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Invalid or expired refresh token"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/auth/refresh [post]
func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}

	userID, err := h.jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	hash := HashToken(req.RefreshToken)
	if _, err := h.svc.RefreshToken(r.Context(), hash); err != nil {
		respondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}
	if err := h.svc.RevokeRefreshToken(r.Context(), hash); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to revoke old refresh token")
		return
	}

	token, err := h.jwt.GenerateToken(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	refreshToken, refreshHash, err := h.jwt.GenerateRefreshToken(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}
	if err := h.svc.StoreRefreshToken(r.Context(), userID, refreshHash, time.Now().Add(168*time.Hour)); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store refresh token")
		return
	}

	respondJSON(w, http.StatusOK, authResponse{
		User:         userResponse{ID: userID},
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// Logout revokes a refresh token
// @Summary Logout
// @Description Revoke a refresh token, preventing further use
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "Refresh token to revoke"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/auth/logout [post]
// @Security BearerAuth
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}

	hash := HashToken(req.RefreshToken)
	if err := h.svc.RevokeRefreshToken(r.Context(), hash); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to revoke refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type userResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type authResponse struct {
	User         userResponse `json:"user"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
}

func toUserResponse(u domain.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
