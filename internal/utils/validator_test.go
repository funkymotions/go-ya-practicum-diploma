package utils

import (
	"testing"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/dto"
	"github.com/stretchr/testify/suite"
)

type ValidatorTestSuite struct {
	suite.Suite
}

func TestValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}

func (s *ValidatorTestSuite) TestValidateJSONBody() {
	type testCase struct {
		name     string
		input    string
		expected *dto.RegisterUserRequest
		wantErr  bool
	}
	testCases := []testCase{
		{
			name:    "valid request",
			input:   `{"login":"user1","password":"pass123"}`,
			wantErr: false,
			expected: &dto.RegisterUserRequest{
				Login:    "user1",
				Password: "pass123",
			},
		},
		{
			name:     "missing login",
			input:    `{"password":"pass123"}`,
			wantErr:  true,
			expected: nil,
		},
		{
			name:     "invalid json",
			input:    `oops...`,
			wantErr:  true,
			expected: nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			actual, err := ValidateJSONBody[dto.RegisterUserRequest](tc.input)
			if tc.wantErr {
				s.Assert().Error(err)
				s.Assert().Nil(actual)
			} else {
				s.Assert().NoError(err)
				s.Assert().Equal(tc.expected, actual)
			}
		})
	}
}
