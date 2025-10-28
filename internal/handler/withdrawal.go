package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/interfaces"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/middleware"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
	"github.com/go-chi/chi/v5"
)

type withdrawalHandler struct {
	withdrawalService interfaces.WithdrawalService
}

func NewWithdrawalHandler(ws interfaces.WithdrawalService) *withdrawalHandler {
	return &withdrawalHandler{
		withdrawalService: ws,
	}
}

func (h *withdrawalHandler) RegisterRoutes(r chi.Router, authMiddleware *middleware.AuthMiddleware) {
	r.
		With(authMiddleware.IsUserAuthenticated).
		Get("/api/user/withdrawals", h.GetUserWithdrawals)
}

func (h *withdrawalHandler) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.RetrieveContextUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	withdrawals, err := h.withdrawalService.GetUserWithdrawals(userID)
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
	responseData := make([]dto.WithdrawalResponse, 0)
	for _, withdrawal := range *withdrawals {
		responseData = append(responseData, dto.WithdrawalResponse{
			Order:       withdrawal.OrderID,
			Sum:         withdrawal.Amount,
			ProcessedAt: withdrawal.ProcessedAt.Format(time.RFC3339),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
