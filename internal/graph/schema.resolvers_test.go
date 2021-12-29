package graph

import (
	"context"
	"testing"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/graph/graphmodel"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMutationResolver_CreateFriend(t *testing.T) {
	tcs := map[string]struct {
		input     graphmodel.Friends
		expResult *graphmodel.IsSuccess
		expError  error
		mockErr   error
	}{
		"success with an input": {
			input: graphmodel.Friends{
				Friends: []string{"andy@example.com", "john@example.com"},
			},
			expResult: &graphmodel.IsSuccess{
				Success: true,
			},
		},
	}
	for desc, testCase := range tcs {
		t.Run(desc, func(t *testing.T) {
			//Given
			ctx := context.Background()
			var mockService SpecService
			mockService.ExpectedCalls = []*mock.Call{
				mockService.On("CreateFriend", mock.Anything, mock.Anything, mock.Anything).Return(testCase.mockErr),
			}

			r := Resolver{
				Service: mockService,
			}
			mut := r.Mutation()

			//When
			result, err := mut.CreateFriend(ctx, testCase.input)

			//Then
			if testCase.expError != nil {
				require.EqualError(t, err, testCase.expError.Error())
			} else {
				require.Equal(t, testCase.expResult, result)
			}
		})
	}
}
