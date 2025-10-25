package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/middleware"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
	"github.com/go-chi/chi/v5"
)

type orderService interface {
	RegisterOrder(*dto.RegisterOrderRequest) error
	GetUserOrders(userID uint) (*[]model.Order, error)
}

type orderHandler struct {
	orderService orderService
}

func NewOrderHandler(os orderService) *orderHandler {
	return &orderHandler{
		orderService: os,
	}
}

func (h *orderHandler) RegisterRoutes(
	r chi.Router,
	authMiddleware *middleware.AuthMiddleware,
) {
	r.
		With(authMiddleware.IsUserAuthenticated).
		Get("/api/user/orders", h.GetUserOrders)
	r.
		With(authMiddleware.IsUserAuthenticated).
		Post("/api/user/orders", h.RegisterOrder)
}

func (h *orderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, ok := utils.RetrieveContextUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	orders, err := h.orderService.GetUserOrders(uint(userID))
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
	var responseData []dto.GetUserOrderResponse
	for _, order := range *orders {
		responseData = append(responseData, dto.GetUserOrderResponse{
			Number:     order.ID,
			Status:     order.OrderStatus,
			UploadedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Accrual:    order.Accrual,
		})
	}
	response, err := json.Marshal(responseData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *orderHandler) RegisterOrder(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckRequestContentType(r, utils.ContentTypePlain) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, ok := utils.RetrieveContextUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderID := string(body)
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req := &dto.RegisterOrderRequest{
		UserID:  uint(userID),
		OrderID: orderID,
	}
	err = h.orderService.RegisterOrder(req)
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
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Order registered successfully"))
}
