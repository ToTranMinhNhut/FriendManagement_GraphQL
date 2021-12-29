package repository

import (
	"context"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Insert a new record into friends table
func (_self DBRepo) CreateFriend(ctx context.Context, userId int, friendId int) error {
	friend := models.Friend{
		UserID:   userId,
		FriendID: friendId,
	}
	return friend.Insert(ctx, _self.Db, boil.Infer())
}

// Get friendship slice from friends table by user id
func (_self DBRepo) GetFriendsByID(ctx context.Context, userId int) (models.FriendSlice, error) {
	return models.Friends(
		qm.Select(models.FriendColumns.UserID, models.FriendColumns.FriendID),
		qm.Where("user_id = ?", userId), qm.Or("friend_id = ?", userId),
	).All(ctx, _self.Db)
}

// Get blocked user relationship slice from user_blocks table by user id
func (_self DBRepo) GetUserBlocksByID(ctx context.Context, userId int) (models.UserBlockSlice, error) {
	return models.UserBlocks(
		qm.Select(models.UserBlockColumns.RequestorID, models.UserBlockColumns.TargetID),
		qm.Where("requestor_id = ?", userId), qm.Or("target_id = ?", userId),
	).All(ctx, _self.Db)
}

// Insert a new record into subscriptions table
func (_self DBRepo) CreateSubscription(ctx context.Context, requestorId int, targetId int) error {
	subscription := models.Subscription{
		SubscriptionRequestorID: requestorId,
		SubscriptionTargetID:    targetId,
	}
	return subscription.Insert(ctx, _self.Db, boil.Infer())
}

// Get users slice (who are not blocked by sender) by user id
func (_self DBRepo) GetRecipientEmails(ctx context.Context, senderId int) (models.UserSlice, error) {
	query := `SELECT DISTINCT val.email FROM (
	        SELECT u.id, u.email
	        FROM users u JOIN friends f ON (u.id = f.user_id OR u.id = f.friend_id)
	        WHERE u.id <> $1 AND (f.user_id = $1 OR f.friend_id = $1)
	        UNION
	        SELECT u.id, u.email
	        FROM subscriptions s JOIN users u ON s.subscription_requestor_id = u.id
	        WHERE u.id <> $1 AND s.subscription_target_id = $1
	    ) AS val
	    WHERE NOT EXISTS(
	        SELECT 1 FROM user_blocks b
	        WHERE (b.requestor_id = val.id AND b.target_id = $1) OR (b.target_id = val.id AND b.requestor_id = $1)
	)`

	nonBlockUsers := models.UserSlice{} //make([]models.User, 0)
	err := queries.Raw(query, senderId).Bind(ctx, _self.Db, &nonBlockUsers)
	if err != nil {
		return nil, err
	}

	return nonBlockUsers, nil
}

// Insert a blocking relationship of users into user_blocks table
func (_self DBRepo) CreateUserBlock(ctx context.Context, requestorId int, targetId int) error {
	userBlock := models.UserBlock{
		RequestorID: requestorId,
		TargetID:    targetId,
	}
	return userBlock.Insert(ctx, _self.Db, boil.Infer())
}

// Verify a existing friendship
func (_self DBRepo) IsExistedFriend(ctx context.Context, userId int, friendId int) (bool, error) {
	return models.Friends(
		qm.WhereIn("user_id in ?", userId, friendId),
		qm.AndIn("friend_id in ?", userId, friendId)).
		Exists(ctx, _self.Db)
}

// Verify a blocking relationship of users
func (_self DBRepo) IsBlockedUser(ctx context.Context, userId int, friendId int) (bool, error) {
	return models.UserBlocks(
		qm.WhereIn("requestor_id in ?", userId, friendId),
		qm.AndIn("target_id in ?", userId, friendId)).
		Exists(ctx, _self.Db)
}

// Verify a subscription relationship of users
func (_self DBRepo) IsSubscribedUser(ctx context.Context, requestorId int, targetId int) (bool, error) {
	return models.Subscriptions(
		qm.WhereIn("subscription_requestor_id in ?", requestorId, targetId),
		qm.AndIn("subscription_target_id in ?", requestorId, targetId)).
		Exists(ctx, _self.Db)
}

// Get a user id from users table by email
func (_self DBRepo) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	var userId int
	user, err := models.Users(qm.Select(models.UserColumns.ID), qm.Where("email = ?", email)).One(ctx, _self.Db)
	if err != nil {
		return userId, err
	}
	return user.ID, nil
}

// Get list of emails by list of corresponding ids from users table
func (_self DBRepo) GetEmailsByUserIDs(ctx context.Context, userIDs []int) ([]string, error) {
	if len(userIDs) == 0 {
		return []string{}, nil
	}

	users, err := models.Users(qm.Select(models.UserColumns.Email), models.UserWhere.ID.IN(userIDs)).All(ctx, _self.Db)
	if err != nil {
		return nil, err
	}

	emails := make([]string, len(users))
	for i, user := range users {
		emails[i] = user.Email
	}
	return emails, nil
}

// Get all users from users table
func (_self DBRepo) GetUsers(ctx context.Context) (models.UserSlice, error) {
	return models.Users().All(ctx, _self.Db)
}
