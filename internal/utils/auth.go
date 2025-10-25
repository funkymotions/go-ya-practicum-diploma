package utils

import (
	"context"
	"strconv"
)

func RetrieveContextUserID(ctx context.Context) (uint, bool) {
	userIDVal := ctx.Value("userID")
	if userIDVal == nil {
		return 0, false
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		return 0, false
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, false
	}
	return uint(userID), true
}
