package http

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret           []byte
	accessExpiresIn  time.Duration
	refreshExpiresIn time.Duration
}

func NewJWTService(secret string, accessExpiresIn, refreshExpiresIn time.Duration) *JWTService {
	return &JWTService{
		secret:           []byte(secret),
		accessExpiresIn:  accessExpiresIn,
		refreshExpiresIn: refreshExpiresIn,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *JWTService) GenerateToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiresIn)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token inválido: método de assinatura inesperado: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("token inválido: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("token inválido")
	}
	return claims.UserID, nil
}

func (s *JWTService) GenerateRefreshToken(userID string) (token string, hash string, err error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiresIn)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString(s.secret)
	if err != nil {
		return "", "", fmt.Errorf("falha ao enviar refresh token: %w", err)
	}
	h := sha256.Sum256([]byte(token))
	hash = fmt.Sprintf("%x", h)
	return token, hash, nil
}

func (s *JWTService) ValidateRefreshToken(tokenStr string) (string, error) {
	return s.ValidateToken(tokenStr)
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

func setRefreshTokenCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/api/auth",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
		Expires:  expires,
	})
}

func clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
		MaxAge:   -1,
	})
}

func (s *JWTService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "header de autorização ausente")
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			respondError(w, http.StatusUnauthorized, "formato de autorização inválido, use Bearer <token>")
			return
		}
		userID, err := s.ValidateToken(parts[1])
		if err != nil {
			respondError(w, http.StatusUnauthorized, "token inválido ou expirado")
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
