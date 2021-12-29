package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServices_GetUsers(t *testing.T) {
	tcs := map[string]struct {
		expResult []string
		expError  error
		mockUsers models.UserSlice
	}{
		"successfully get all user": {
			mockUsers: models.UserSlice{
				&models.User{Name: "john", Email: "john@example.com"},
				&models.User{Name: "andy", Email: "andy@example.com"},
			},
			expResult: []string{"john@example.com", "andy@example.com"},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			var mockRepo SpecRepo
			mockRepo.ExpectedCalls = []*mock.Call{
				mockRepo.On("GetUsers", mock.Anything).Return(tc.mockUsers, tc.expError),
			}

			friendService := NewFriendService(mockRepo)
			result, err := friendService.GetUsers(ctx)
			require.NoError(t, err)
			require.Equal(t, len(tc.expResult), len(result))
			for i, ss := range tc.expResult {
				require.Equal(t, ss, result[i])
			}
		})
	}
}

func TestServices_CreateFriends(t *testing.T) {
	type mockGetUserID struct {
		result int
		err    error
	}
	type mockIsExistedFriend struct {
		result bool
		err    error
	}
	type mockIsBlockedUser struct {
		result bool
		err    error
	}

	tcs := map[string]struct {
		userEmail       string
		friendEmail     string
		firstUser       mockGetUserID
		secondUser      mockGetUserID
		isExistedFriend mockIsExistedFriend
		isBlockedUser   mockIsBlockedUser
		expError        error
	}{
		"success with an input": {
			userEmail:   "andy@example.com",
			friendEmail: "john@example.com",
			firstUser: mockGetUserID{
				result: 101,
			},
			secondUser: mockGetUserID{
				result: 100,
			},
			isExistedFriend: mockIsExistedFriend{
				result: false,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
		},
		"failed with an unknow format input of user": {
			userEmail:   "test@example.com",
			friendEmail: "john@example.com",
			firstUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			secondUser: mockGetUserID{
				result: 100,
			},
			isExistedFriend: mockIsExistedFriend{
				result: false,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
			expError: errors.New(`test@example.com is not exists`),
		},
		"failed with an unknow format input of friend": {
			userEmail:   "andy@example.com",
			friendEmail: "test@example.com",
			firstUser: mockGetUserID{
				result: 101,
			},
			secondUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			isExistedFriend: mockIsExistedFriend{
				result: false,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
			expError: errors.New(`test@example.com is not exists`),
		},
		"failed with an friendship is existing": {
			userEmail:   "john@example.com",
			friendEmail: "andy@example.com",
			firstUser: mockGetUserID{
				result: 100,
			},
			secondUser: mockGetUserID{
				result: 101,
			},
			isExistedFriend: mockIsExistedFriend{
				result: true,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
			expError: errors.New(`The friend relationship has been existed`),
		},
		"failed with an blocking relationship is existing": {
			userEmail:   "john@example.com",
			friendEmail: "andy@example.com",
			firstUser: mockGetUserID{
				result: 100,
			},
			secondUser: mockGetUserID{
				result: 101,
			},
			isExistedFriend: mockIsExistedFriend{
				result: false,
			},
			isBlockedUser: mockIsBlockedUser{
				result: true,
			},
			expError: errors.New(`The users have blocked each other`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			var mockRepo SpecRepo
			mockRepo.ExpectedCalls = []*mock.Call{
				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.firstUser.result, tc.firstUser.err).Once(),
				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.secondUser.result, tc.secondUser.err),
				mockRepo.On("IsExistedFriend", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.isExistedFriend.result, tc.isExistedFriend.err),
				mockRepo.On("IsBlockedUser", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.isBlockedUser.result, tc.isBlockedUser.err),
				mockRepo.On("CreateFriend", mock.Anything, mock.Anything, mock.Anything).
					Return(nil),
			}
			friendService := NewFriendService(mockRepo)
			err := friendService.CreateFriend(ctx, tc.userEmail, tc.friendEmail)
			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}

func TestServices_GetFriends(t *testing.T) {
	type mockGetUserID struct {
		result int
		err    error
	}
	type mockGetFriends struct {
		result models.FriendSlice
		err    error
	}
	type mockGetUserBlocks struct {
		result models.UserBlockSlice
		err    error
	}
	type mockGetEmails struct {
		result []string
		err    error
	}

	tcs := map[string]struct {
		userEmail      string
		expResult      []string
		expError       error
		mockUser       mockGetUserID
		mockFriends    mockGetFriends
		mockUserBlocks mockGetUserBlocks
		mockEmails     mockGetEmails
	}{
		"success with an input": {
			userEmail: "andy@example.com",
			expResult: []string{"john@example.com"},
			mockUser: mockGetUserID{
				result: 100,
			},
			mockFriends: mockGetFriends{
				result: models.FriendSlice{
					&models.Friend{UserID: 100, FriendID: 101},
					&models.Friend{UserID: 100, FriendID: 102},
				},
			},
			mockUserBlocks: mockGetUserBlocks{
				result: models.UserBlockSlice{
					&models.UserBlock{RequestorID: 100, TargetID: 102},
				},
			},
			mockEmails: mockGetEmails{
				result: []string{"john@example.com"},
			},
		},
		"failed with an unknow format input": {
			userEmail: "test@example.com",
			mockUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			expError: errors.New(`test@example.com is not exists`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			var mockRepo SpecRepo
			mockRepo.ExpectedCalls = []*mock.Call{
				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.mockUser.result, tc.mockUser.err),
				mockRepo.On("GetFriendsByID", mock.Anything, mock.Anything).
					Return(tc.mockFriends.result, tc.mockFriends.err),
				mockRepo.On("GetUserBlocksByID", mock.Anything, mock.Anything).
					Return(tc.mockUserBlocks.result, tc.mockUserBlocks.err),
				mockRepo.On("GetEmailsByUserIDs", mock.Anything, mock.Anything).
					Return(tc.mockEmails.result, tc.mockEmails.err),
			}
			friendService := NewFriendService(mockRepo)
			result, err := friendService.GetFriends(ctx, tc.userEmail)
			if tc.expError != nil {
				require.EqualError(t, tc.expError, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestServices_GetCommonFriends(t *testing.T) {
	type mockGetUserID struct {
		result int
		err    error
	}
	type mockGetFriends struct {
		result models.FriendSlice
		err    error
	}
	type mockGetUserBlocks struct {
		result models.UserBlockSlice
		err    error
	}
	type mockGetEmails struct {
		result []string
		err    error
	}

	tcs := map[string]struct {
		firstEmail           string
		secondEmail          string
		expResult            []string
		expError             error
		firstMockUser        mockGetUserID
		firstMockFriends     mockGetFriends
		firstMockUserBlocks  mockGetUserBlocks
		firstMockEmails      mockGetEmails
		secondMockUser       mockGetUserID
		secondMockFriends    mockGetFriends
		secondMockUserBlocks mockGetUserBlocks
		secondMockEmails     mockGetEmails
	}{
		"success with an input": {
			firstEmail:  "john@example.com",
			secondEmail: "lisa@example.com",
			expResult:   []string{"common@example.com"},
			firstMockUser: mockGetUserID{
				result: 100,
			},
			firstMockFriends: mockGetFriends{
				result: models.FriendSlice{
					&models.Friend{UserID: 100, FriendID: 101},
					&models.Friend{UserID: 100, FriendID: 102},
				},
			},
			firstMockUserBlocks: mockGetUserBlocks{
				result: models.UserBlockSlice{
					&models.UserBlock{RequestorID: 100, TargetID: 101},
				},
			},
			firstMockEmails: mockGetEmails{
				result: []string{"common@example.com"},
			},
			secondMockUser: mockGetUserID{
				result: 100,
			},
			secondMockFriends: mockGetFriends{
				result: models.FriendSlice{
					&models.Friend{UserID: 103, FriendID: 101},
					&models.Friend{UserID: 103, FriendID: 102},
				},
			},
			secondMockUserBlocks: mockGetUserBlocks{
				result: models.UserBlockSlice{
					&models.UserBlock{RequestorID: 103, TargetID: 100},
				},
			},
			secondMockEmails: mockGetEmails{
				result: []string{"common@example.com"},
			},
		},
		"failed with an unknow format input of first user": {
			firstEmail: "test@example.com",
			firstMockUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			expError: errors.New(`test@example.com is not exists`),
		},
		"failed with an unknow format input of second user": {
			firstEmail:  "john@example.com",
			secondEmail: "test@example.com",
			firstMockUser: mockGetUserID{
				result: 100,
			},
			secondMockUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			expError: errors.New(`test@example.com is not exists`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			var mockRepo SpecRepo
			mockRepo.ExpectedCalls = []*mock.Call{
				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.firstMockUser.result, tc.firstMockUser.err).Once(),
				mockRepo.On("GetFriendsByID", mock.Anything, mock.Anything).
					Return(tc.firstMockFriends.result, tc.firstMockFriends.err).Once(),
				mockRepo.On("GetUserBlocksByID", mock.Anything, mock.Anything).
					Return(tc.firstMockUserBlocks.result, tc.firstMockUserBlocks.err).Once(),
				mockRepo.On("GetEmailsByUserIDs", mock.Anything, mock.Anything).
					Return(tc.firstMockEmails.result, tc.firstMockEmails.err).Once(),

				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.secondMockUser.result, tc.secondMockUser.err),
				mockRepo.On("GetFriendsByID", mock.Anything, mock.Anything).
					Return(tc.secondMockFriends.result, tc.secondMockFriends.err),
				mockRepo.On("GetUserBlocksByID", mock.Anything, mock.Anything).
					Return(tc.secondMockUserBlocks.result, tc.secondMockUserBlocks.err),
				mockRepo.On("GetEmailsByUserIDs", mock.Anything, mock.Anything).
					Return(tc.secondMockEmails.result, tc.secondMockEmails.err),
			}
			friendService := NewFriendService(mockRepo)
			result, err := friendService.GetCommonFriends(ctx, tc.firstEmail, tc.secondEmail)
			if tc.expError != nil {
				require.EqualError(t, tc.expError, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestServices_CreateSubcription(t *testing.T) {
	type mockGetUserID struct {
		result int
		err    error
	}
	type mockIsSubscribedUser struct {
		result bool
		err    error
	}
	type mockIsBlockedUser struct {
		result bool
		err    error
	}

	tcs := map[string]struct {
		requestorEmail   string
		targetEmail      string
		firstUser        mockGetUserID
		secondUser       mockGetUserID
		isSubscribedUser mockIsSubscribedUser
		isBlockedUser    mockIsBlockedUser
		expError         error
	}{
		"success with an input": {
			requestorEmail: "andy@example.com",
			targetEmail:    "john@example.com",
			firstUser: mockGetUserID{
				result: 101,
			},
			secondUser: mockGetUserID{
				result: 100,
			},
			isSubscribedUser: mockIsSubscribedUser{
				result: false,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
		},
		"failed with an unknow format input of requestor": {
			requestorEmail: "test@example.com",
			targetEmail:    "john@example.com",
			firstUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			secondUser: mockGetUserID{
				result: 100,
			},
			isSubscribedUser: mockIsSubscribedUser{
				result: false,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
			expError: errors.New(`test@example.com is not exists`),
		},
		"failed with an unknow format input of target user": {
			requestorEmail: "andy@example.com",
			targetEmail:    "test@example.com",
			firstUser: mockGetUserID{
				result: 101,
			},
			secondUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			isSubscribedUser: mockIsSubscribedUser{
				result: false,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
			expError: errors.New(`test@example.com is not exists`),
		},
		"failed with an subscribed relatinship is existing": {
			requestorEmail: "john@example.com",
			targetEmail:    "andy@example.com",
			firstUser: mockGetUserID{
				result: 100,
			},
			secondUser: mockGetUserID{
				result: 101,
			},
			isSubscribedUser: mockIsSubscribedUser{
				result: true,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
			expError: errors.New(`The users have subscribed each other`),
		},
		"failed with an blocking relationship is existing": {
			requestorEmail: "john@example.com",
			targetEmail:    "andy@example.com",
			firstUser: mockGetUserID{
				result: 100,
			},
			secondUser: mockGetUserID{
				result: 101,
			},
			isSubscribedUser: mockIsSubscribedUser{
				result: false,
			},
			isBlockedUser: mockIsBlockedUser{
				result: true,
			},
			expError: errors.New(`The users have blocked each other`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			var mockRepo SpecRepo
			mockRepo.ExpectedCalls = []*mock.Call{
				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.firstUser.result, tc.firstUser.err).Once(),
				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.secondUser.result, tc.secondUser.err),
				mockRepo.On("IsSubscribedUser", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.isSubscribedUser.result, tc.isSubscribedUser.err),
				mockRepo.On("IsBlockedUser", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.isBlockedUser.result, tc.isBlockedUser.err),
				mockRepo.On("CreateSubscription", mock.Anything, mock.Anything, mock.Anything).
					Return(nil),
			}
			friendService := NewFriendService(mockRepo)
			err := friendService.CreateSubscription(ctx, tc.requestorEmail, tc.targetEmail)
			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}

func TestServices_CreateUserBlocks(t *testing.T) {
	type mockGetUserID struct {
		result int
		err    error
	}
	type mockIsBlockedUser struct {
		result bool
		err    error
	}

	tcs := map[string]struct {
		requestorEmail string
		targetEmail    string
		firstUser      mockGetUserID
		secondUser     mockGetUserID
		isBlockedUser  mockIsBlockedUser
		expError       error
	}{
		"success with an input": {
			requestorEmail: "andy@example.com",
			targetEmail:    "john@example.com",
			firstUser: mockGetUserID{
				result: 101,
			},
			secondUser: mockGetUserID{
				result: 100,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
		},
		"failed with an unknow format input of requestor": {
			requestorEmail: "test@example.com",
			targetEmail:    "john@example.com",
			firstUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			secondUser: mockGetUserID{
				result: 100,
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
			expError: errors.New(`test@example.com is not exists`),
		},
		"failed with an unknow format input of target user": {
			requestorEmail: "andy@example.com",
			targetEmail:    "test@example.com",
			firstUser: mockGetUserID{
				result: 101,
			},
			secondUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			isBlockedUser: mockIsBlockedUser{
				result: false,
			},
			expError: errors.New(`test@example.com is not exists`),
		},
		"failed with an blocking relationship is existing": {
			requestorEmail: "john@example.com",
			targetEmail:    "andy@example.com",
			firstUser: mockGetUserID{
				result: 100,
			},
			secondUser: mockGetUserID{
				result: 101,
			},
			isBlockedUser: mockIsBlockedUser{
				result: true,
			},
			expError: errors.New(`The users have blocked each other`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			var mockRepo SpecRepo
			mockRepo.ExpectedCalls = []*mock.Call{
				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.firstUser.result, tc.firstUser.err).Once(),

				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.secondUser.result, tc.secondUser.err),

				mockRepo.On("IsBlockedUser", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.isBlockedUser.result, tc.isBlockedUser.err),

				mockRepo.On("CreateUserBlock", mock.Anything, mock.Anything, mock.Anything).
					Return(nil),
			}
			friendService := NewFriendService(mockRepo)
			err := friendService.CreateUserBlock(ctx, tc.requestorEmail, tc.targetEmail)
			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}

func TestServices_GetRecipientEmails(t *testing.T) {
	type mockGetUserID struct {
		result int
		err    error
	}
	type mockGetFriends struct {
		result models.FriendSlice
		err    error
	}
	type mockGetRecipients struct {
		result models.UserSlice
		err    error
	}

	tcs := map[string]struct {
		userEmail      string
		text           string
		expResult      []string
		expError       error
		mockUser       mockGetUserID
		mockRecipients mockGetRecipients
	}{
		"success with an input": {
			userEmail: "andy@example.com",
			text:      "hello! kate@example.com",
			expResult: []string{"john@example.com", "kate@example.com"},
			mockUser: mockGetUserID{
				result: 100,
			},
			mockRecipients: mockGetRecipients{
				result: models.UserSlice{
					&models.User{Name: "john", Email: "john@example.com"},
				},
			},
		},
		"failed with an unknow format input": {
			userEmail: "test@example.com",
			mockUser: mockGetUserID{
				err: errors.New(`test@example.com is not exists`),
			},
			expError: errors.New(`test@example.com is not exists`),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			var mockRepo SpecRepo
			mockRepo.ExpectedCalls = []*mock.Call{
				mockRepo.On("GetUserIDByEmail", mock.Anything, mock.Anything).
					Return(tc.mockUser.result, tc.mockUser.err),

				mockRepo.On("GetRecipientEmails", mock.Anything, mock.Anything).
					Return(tc.mockRecipients.result, tc.mockRecipients.err),
			}
			friendService := NewFriendService(mockRepo)
			result, err := friendService.GetRecipientEmails(ctx, tc.userEmail, tc.text)
			if tc.expError != nil {
				require.EqualError(t, tc.expError, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}
