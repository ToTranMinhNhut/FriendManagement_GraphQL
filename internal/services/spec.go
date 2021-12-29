package services

import (
	"context"
)

// SpecRepo is the interface for repository methods
type SpecService interface {
	CreateFriend(ctx context.Context, userEmail string, friendEmail string) error
	GetFriends(ctx context.Context, userEmail string) ([]string, error)
	GetCommonFriends(ctx context.Context, firstUserEmail string, secondUserEmail string) ([]string, error)
	CreateSubscription(ctx context.Context, requestorEmail string, targetEmail string) error
	CreateUserBlock(ctx context.Context, requestorEmail string, targetEmail string) error
	GetRecipientEmails(ctx context.Context, senderEmail string, text string) ([]string, error)
	GetUsers(ctx context.Context) ([]string, error)
}
