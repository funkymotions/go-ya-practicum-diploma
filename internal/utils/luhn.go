package utils

import (
	"strconv"
	"strings"
)

func IsValidLuhn(number string) bool {
	digits := make([]int64, len(number))
	rawDigits := strings.Split(number, "")
	for i, d := range rawDigits {
		digit, _ := strconv.ParseInt(d, 10, 64)
		digits[i] = digit
	}
	sum := int64(0)
	double := false
	for i := len(digits) - 1; i >= 0; i-- {
		d := digits[i]
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}
