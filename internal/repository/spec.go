package repository

import (
	"context"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/models"
)

// SpecRepo is the interface for repository methods
type SpecRepo interface {
	CreateFriend(ctx context.Context, userId int, friendId int) error
	GetFriendsByID(ctx context.Context, userId int) (models.FriendSlice, error)
	GetUserBlocksByID(ctx context.Context, userId int) (models.UserBlockSlice, error)
	CreateSubscription(ctx context.Context, requestorId int, targetId int) error
	GetRecipientEmails(ctx context.Context, senderId int) (models.UserSlice, error)
	CreateUserBlock(ctx context.Context, requestorId int, targetId int) error
	IsExistedFriend(ctx context.Context, userId int, friendId int) (bool, error)
	IsBlockedUser(ctx context.Context, userId int, friendId int) (bool, error)
	IsSubscribedUser(ctx context.Context, requestorId int, targetId int) (bool, error)
	GetUserIDByEmail(ctx context.Context, email string) (int, error)
	GetEmailsByUserIDs(ctx context.Context, userIDs []int) ([]string, error)
	GetUsers(ctx context.Context) (models.UserSlice, error)
}
