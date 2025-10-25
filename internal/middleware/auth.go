package middleware

import (
	"context"
	"net/http"
	"strconv"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type UserIDValue string

type userRepository interface {
	FindOneByID(id uint) (*model.User, error)
}

type AuthMiddleware struct {
	userRepository userRepository
}

func NewAuthUserMiddleware(r userRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userRepository: r,
	}
}

func (m *AuthMiddleware) IsUserAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("gophermart_session")
		if err != nil || cookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID := cookie.Value
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, err = m.userRepository.FindOneByID(uint(userIDInt))
		if err == apperrors.ErrUserNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
