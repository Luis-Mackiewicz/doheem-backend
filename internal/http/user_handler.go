package http

import (
	"errors"
	"net/http"
	"time"

	"doheem-backend/internal/user"
)

type UserHandler struct {
	svc *user.UserService
	jwt *JWTService
}

func NewUserHandler(svc *user.UserService, jwt *JWTService) *UserHandler {
	return &UserHandler{svc: svc, jwt: jwt}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string  `json:"name"       validate:"required"`
		Email     string  `json:"email"      validate:"required,email"`
		Password  string  `json:"password"   validate:"required,min=6"`
		Phone     string  `json:"phone"      validate:"required,phone_br"`
		Document  string  `json:"document"   validate:"required,document"`
		BirthDate string  `json:"birth_date" validate:"required"`
		Cep       string  `json:"cep"        validate:"required,cep_br"`
		AvatarURL *string `json:"avatar_url,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "formato de birth_date inválido, use AAAA-MM-DD")
		return
	}
	phone := onlyDigits(req.Phone)
	document := onlyDigits(req.Document)
	cep := onlyDigits(req.Cep)
	created, err := h.svc.Register(r.Context(), user.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Phone:        &phone,
		Document:     &document,
		BirthDate:    &birthDate,
		Cep:          &cep,
		AvatarURL:    req.AvatarURL,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}
	token, err := h.jwt.GenerateToken(created.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao gerar token de autenticação")
		return
	}
	refreshToken, refreshHash, err := h.jwt.GenerateRefreshToken(created.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao gerar token de atualização")
		return
	}
	if err := h.svc.StoreRefreshToken(r.Context(), created.ID, refreshHash, time.Now().Add(168*time.Hour)); err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao armazenar token de atualização")
		return
	}
	setRefreshTokenCookie(w, refreshToken, time.Now().Add(168*time.Hour))
	respondJSON(w, http.StatusCreated, authResponse{
		User:  toUserResponse(created),
		Token: token,
	})
}

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
		respondError(w, http.StatusUnauthorized, "e-mail ou senha inválidos")
		return
	}
	token, err := h.jwt.GenerateToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao gerar token de autenticação")
		return
	}
	refreshToken, refreshHash, err := h.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao gerar token de atualização")
		return
	}
	if err := h.svc.StoreRefreshToken(r.Context(), user.ID, refreshHash, time.Now().Add(168*time.Hour)); err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao armazenar token de atualização")
		return
	}
	setRefreshTokenCookie(w, refreshToken, time.Now().Add(168*time.Hour))
	respondJSON(w, http.StatusOK, authResponse{
		User:  toUserResponse(user),
		Token: token,
	})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	user, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	var req struct {
		Name      *string `json:"name,omitempty"`
		Email     *string `json:"email,omitempty"     validate:"omitempty,email"`
		Phone     *string `json:"phone,omitempty"     validate:"omitempty,phone_br"`
		Document  *string `json:"document,omitempty" validate:"omitempty,document"`
		BirthDate *string `json:"birth_date,omitempty"`
		Cep       *string `json:"cep,omitempty"     validate:"omitempty,cep_br"`
		AvatarURL *string `json:"avatar_url,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	var birthDate *time.Time
	if req.BirthDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "formato de birth_date inválido, use AAAA-MM-DD")
			return
		}
		birthDate = &parsed
	}

	var phone *string
	if req.Phone != nil {
		v := onlyDigits(*req.Phone)
		phone = &v
	}
	var document *string
	if req.Document != nil {
		v := onlyDigits(*req.Document)
		document = &v
	}
	var cep *string
	if req.Cep != nil {
		v := onlyDigits(*req.Cep)
		cep = &v
	}

	updated, err := h.svc.Update(r.Context(), userID, user.UpdateUserParams{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     phone,
		Document:  document,
		BirthDate: birthDate,
		Cep:       cep,
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toUserResponse(updated))
}

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
		if errors.Is(err, user.ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "senha antiga inválida")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := ""
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		refreshToken = cookie.Value
	}
	if refreshToken == "" {
		respondError(w, http.StatusBadRequest, "token de atualização é obrigatório")
		return
	}

	userID, err := h.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "token de atualização inválido ou expirado")
		return
	}

	hash := HashToken(refreshToken)
	if _, err := h.svc.RefreshToken(r.Context(), hash); err != nil {
		respondError(w, http.StatusUnauthorized, "token de atualização inválido ou expirado")
		return
	}
	if err := h.svc.RevokeRefreshToken(r.Context(), hash); err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao revogar token de atualização antigo")
		return
	}

	user, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "usuário não encontrado")
		return
	}

	token, err := h.jwt.GenerateToken(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao gerar token de autenticação")
		return
	}
	newRefreshToken, refreshHash, err := h.jwt.GenerateRefreshToken(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao gerar token de atualização")
		return
	}
	if err := h.svc.StoreRefreshToken(r.Context(), userID, refreshHash, time.Now().Add(168*time.Hour)); err != nil {
		respondError(w, http.StatusInternalServerError, "falha ao armazenar token de atualização")
		return
	}

	setRefreshTokenCookie(w, newRefreshToken, time.Now().Add(168*time.Hour))
	respondJSON(w, http.StatusOK, authResponse{
		User:  toUserResponse(user),
		Token: token,
	})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken := ""
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		refreshToken = cookie.Value
	}
	if refreshToken == "" {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if errs := decodeAndValidate(r, &req); errs == nil && req.RefreshToken != "" {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken != "" {
		hash := HashToken(refreshToken)
		h.svc.RevokeRefreshToken(r.Context(), hash)
	}

	clearRefreshTokenCookie(w)
	respondJSON(w, http.StatusOK, map[string]string{"message": "logout bem-sucedido"})
}

type userResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Document  string  `json:"document"`
	BirthDate string  `json:"birth_date"`
	Cep       string  `json:"cep"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type authResponse struct {
	User  userResponse `json:"user"`
	Token string       `json:"token"`
}

func toUserResponse(u user.User) userResponse {
	birthDate := ""
	if u.BirthDate != nil {
		birthDate = u.BirthDate.Format("2006-01-02")
	}
	return userResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     maskPhone(strVal(u.Phone)),
		Document:  maskDocument(strVal(u.Document)),
		BirthDate: birthDate,
		Cep:       maskCEP(strVal(u.Cep)),
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
