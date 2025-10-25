package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
	"github.com/go-chi/chi/v5"
)

type userService interface {
	Register(*dto.RegisterUserRequest) (*model.User, error)
	Login(*dto.LoginUserRequest) (*model.User, error)
}

type userHandler struct {
	userService userService
}

func NewUserHandler(us userService) *userHandler {
	return &userHandler{
		userService: us,
	}
}

func (h *userHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Login)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckRequestContentType(r, utils.ContentTypeJSON) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req, err := utils.ValidateJSONBody[dto.LoginUserRequest](string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.userService.Login(req)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			w.WriteHeader(appErr.StatusCode)
			w.Write([]byte(appErr.Message))
			return
		}
		fmt.Printf("unexpected error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "gophermart_session",
		Value:    fmt.Sprintf("%d", user.ID),
		Path:     "/",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckRequestContentType(r, utils.ContentTypeJSON) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req, err := utils.ValidateJSONBody[dto.RegisterUserRequest](string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.userService.Register(req)
	if err != nil {
		// aknowledge application errors
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			w.WriteHeader(appErr.StatusCode)
			w.Write([]byte(appErr.Message))
			return
		}
		fmt.Printf("unexpected error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "gophermart_session",
		Value:    fmt.Sprintf("%d", user.ID),
		Path:     "/",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
