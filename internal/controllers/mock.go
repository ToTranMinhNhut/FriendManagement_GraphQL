package controllers

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type SpecService struct {
	mock.Mock
}

func (m SpecService) CreateFriend(ctx context.Context, userEmail string, friendEmail string) error {
	args := m.Called(ctx, userEmail, friendEmail)
	var r error
	if args.Get(0) != nil {
		r = args.Get(0).(error)
	}
	return r
}

func (m SpecService) GetFriends(ctx context.Context, userEmail string) ([]string, error) {
	args := m.Called(ctx, userEmail)
	r1 := args.Get(0).([]string)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecService) GetCommonFriends(ctx context.Context, firstUserEmail string, secondUserEmail string) ([]string, error) {
	args := m.Called(ctx, firstUserEmail, secondUserEmail)
	r1 := args.Get(0).([]string)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecService) CreateSubscription(ctx context.Context, requestorEmail string, targetEmail string) error {
	args := m.Called(ctx, requestorEmail, targetEmail)
	var r error
	if args.Get(0) != nil {
		r = args.Get(0).(error)
	}
	return r
}

func (m SpecService) CreateUserBlock(ctx context.Context, requestorEmail string, targetEmail string) error {
	args := m.Called(ctx, requestorEmail, targetEmail)
	var r error
	if args.Get(0) != nil {
		r = args.Get(0).(error)
	}
	return r
}

func (m SpecService) GetRecipientEmails(ctx context.Context, senderEmail string, text string) ([]string, error) {
	args := m.Called(ctx, senderEmail, text)
	r1 := args.Get(0).([]string)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}

func (m SpecService) GetUsers(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	r1 := args.Get(0).([]string)

	var r2 error
	if args.Get(1) != nil {
		r2 = args.Get(1).(error)
	}
	return r1, r2
}
