package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) TestUserContextRetrieval() {
	type testCase struct {
		name       string
		userID     interface{}
		wantError  bool
		wantUserID uint
	}
	testCases := []testCase{
		{
			name:      "Invalid user ID",
			userID:    1,
			wantError: true,
		},
		{
			name:       "Valid user ID",
			userID:     "1",
			wantError:  false,
			wantUserID: 1,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "userID", tc.userID)
			userID, ok := RetrieveContextUserID(ctx)
			if tc.wantError {
				s.Assert().False(ok, "Expected error but got valid user ID")
			} else {
				s.Assert().True(ok, "Expected valid user ID but got error")
				s.Assert().Equal(tc.wantUserID, userID, "User ID does not match expected value")
			}
		})
	}
}
