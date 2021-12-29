package services

import (
	"context"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/models"
	"github.com/stretchr/testify/mock"
)

type SpecRepo struct {
	mock.Mock
}

func (m SpecRepo) GetUsers(ctx context.Context) (models.UserSlice, error) {
	args := m.Called(ctx)
	r1 := args.Get(0).(models.UserSlice)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecRepo) CreateFriend(ctx context.Context, userId int, friendId int) error {
	args := m.Called(ctx, userId, friendId)
	var r error
	if args.Get(0) != nil {
		r = args.Get(0).(error)
	}
	return r
}

func (m SpecRepo) GetFriendsByID(ctx context.Context, userId int) (models.FriendSlice, error) {
	args := m.Called(userId)
	t := args.Get(0).(models.FriendSlice)
	return t, args.Error(1)
}

func (m SpecRepo) GetUserBlocksByID(ctx context.Context, userId int) (models.UserBlockSlice, error) {
	args := m.Called(userId)
	r1 := args.Get(0).(models.UserBlockSlice)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecRepo) CreateSubscription(ctx context.Context, requestorId int, targetId int) error {
	args := m.Called(ctx, requestorId, targetId)
	var r error
	if args.Get(0) != nil {
		r = args.Get(0).(error)
	}
	return r
}

func (m SpecRepo) GetRecipientEmails(ctx context.Context, senderId int) (models.UserSlice, error) {
	args := m.Called(ctx, senderId)
	r1 := args.Get(0).(models.UserSlice)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecRepo) CreateUserBlock(ctx context.Context, requestorId int, targetId int) error {
	args := m.Called(ctx, requestorId, targetId)
	var r error
	if args.Get(0) != nil {
		r = args.Get(0).(error)
	}
	return r
}

func (m SpecRepo) IsExistedFriend(ctx context.Context, userId int, friendId int) (bool, error) {
	args := m.Called(ctx, userId, friendId)
	r1 := args.Get(0).(bool)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecRepo) IsBlockedUser(ctx context.Context, userId int, friendId int) (bool, error) {
	args := m.Called(ctx, userId, friendId)
	r1 := args.Get(0).(bool)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecRepo) IsSubscribedUser(ctx context.Context, requestorId int, targetId int) (bool, error) {
	args := m.Called(ctx, requestorId, targetId)
	r1 := args.Get(0).(bool)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecRepo) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	args := m.Called(email)

	t := args.Get(0).(int)
	return t, args.Error(1)
}

func (m SpecRepo) GetEmailsByUserIDs(ctx context.Context, userIDs []int) ([]string, error) {
	args := m.Called(userIDs)
	r1 := args.Get(0).([]string)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}
