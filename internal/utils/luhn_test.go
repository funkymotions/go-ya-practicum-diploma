package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type LuhnTestSuite struct {
	suite.Suite
}

func TestLuhn(t *testing.T) {
	suite.Run(t, new(LuhnTestSuite))
}

func (s *LuhnTestSuite) TestIsValidLuhn() {
	type testCase struct {
		number  string
		isValid bool
	}
	testCases := []testCase{
		{"79927398713", true},
		{"1234567812345670", true},
		{"1234567812345678", false},
		{"49927398716", true},
	}
	for _, tc := range testCases {
		s.Run(tc.number, func() {
			result := IsValidLuhn(tc.number)
			s.Assert().Equal(tc.isValid, result, "Luhn validation failed for number: %s", tc.number)
		})
	}
}
