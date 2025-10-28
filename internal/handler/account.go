package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/interfaces"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/middleware"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
	"github.com/go-chi/chi/v5"
)

type accountHandler struct {
	accountService interfaces.AccountService
}

func NewAccountHandler(as interfaces.AccountService) *accountHandler {
	return &accountHandler{
		accountService: as,
	}
}

func (h *accountHandler) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	r.
		With(authMiddleware.IsUserAuthenticated).
		Get("/api/user/balance", h.GetAccountBalance)
	r.
		With(authMiddleware.IsUserAuthenticated).
		Post("/api/user/balance/withdraw", h.WithdrawAccount)
}

func (h *accountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := utils.RetrieveContextUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, withdrawn, err := h.accountService.CalculateAccountBalance(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var responseData dto.AccountBalanceResponse
	responseData.Current = balance
	responseData.Withdrawn = withdrawn
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *accountHandler) WithdrawAccount(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckRequestContentType(r, utils.ContentTypeJSON) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	userID, ok := utils.RetrieveContextUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := utils.ValidateJSONBody[dto.AccountBalanceWithdrawRequest](string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.accountService.WithdrawFromAccount(userID, data)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			w.WriteHeader(appErr.StatusCode)
			w.Write([]byte(appErr.Message))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
