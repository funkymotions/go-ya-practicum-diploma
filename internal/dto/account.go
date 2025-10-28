package dto

type AccountBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type AccountBalanceWithdrawRequest struct {
	Order string  `json:"order" validate:"required,numeric"`
	Sum   float64 `json:"sum" validate:"required,numeric,gt=0"`
}
