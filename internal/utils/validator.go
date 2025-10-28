package utils

import (
	"encoding/json"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidatedJSONRequest interface {
	dto.RegisterUserRequest | dto.LoginUserRequest | dto.AccountBalanceWithdrawRequest
}

func ValidateJSONBody[T ValidatedJSONRequest](input string) (*T, error) {
	var requestBody T
	if err := json.Unmarshal([]byte(input), &requestBody); err != nil {
		return nil, err
	}
	if err := validate.Struct(requestBody); err != nil {
		return nil, err
	}
	return &requestBody, nil
}
