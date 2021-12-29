package services

import (
	"context"
	"net/http"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/errs"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/utils"
)

// Get all emails of users from repository
func (_self FriendService) GetUsers(ctx context.Context) ([]string, error) {
	users, err := _self.Repo.GetUsers(ctx)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}

	emails := []string{}
	for _, user := range users {
		emails = append(emails, user.Email)
	}

	return emails, nil
}

// Create a new friendship by user id and friend id
func (_self FriendService) CreateFriend(ctx context.Context, userEmail string, friendEmail string) error {
	// Get user id and friend id from repository
	userId, err := _self.Repo.GetUserIDByEmail(ctx, userEmail)
	if err != nil {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: userEmail + " is not exists"}
	}
	friendId, err := _self.Repo.GetUserIDByEmail(ctx, friendEmail)
	if err != nil {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: friendEmail + " is not exists"}
	}

	// Check friend relationship is exists
	isExisted, err := _self.Repo.IsExistedFriend(ctx, userId, friendId)
	if err != nil {
		return &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}
	if isExisted {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: errs.MsgExistedFriendship}
	}

	// Check blocking between 2 emails
	isBlocked, err := _self.Repo.IsBlockedUser(ctx, userId, friendId)
	if err != nil {
		return &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}
	if isBlocked {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: errs.MsgExistedBlockedUser}
	}

	if err := _self.Repo.CreateFriend(ctx, userId, friendId); err != nil {
		return &errs.FriendError{Code: http.StatusInternalServerError, Description: errs.MsgCreatedFriendship}
	}

	return nil
}

// Get email of all friends from a user
func (_self FriendService) GetFriends(ctx context.Context, userEmail string) ([]string, error) {
	// Get user id from an email
	userId, err := _self.Repo.GetUserIDByEmail(ctx, userEmail)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusBadGateway, Description: userEmail + " is not exists"}
	}

	// Get friends available
	friendEmails, err := _self.getFriendEmailsWithoutBlocking(ctx, userId)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}

	return friendEmails, nil
}

// Get emails of common friends by first user and second user
func (_self FriendService) GetCommonFriends(ctx context.Context, firstUserEmail string, secondUserEmail string) ([]string, error) {
	// Get user id and friend id from repository
	firstUserID, err := _self.Repo.GetUserIDByEmail(ctx, firstUserEmail)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusBadGateway, Description: firstUserEmail + " is not exists"}
	}
	secondUserID, err := _self.Repo.GetUserIDByEmail(ctx, secondUserEmail)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusBadGateway, Description: secondUserEmail + " is not exists"}
	}

	// Get friends of first user and second user
	firstFriendEmails, err := _self.getFriendEmailsWithoutBlocking(ctx, firstUserID)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}
	secondFriendEmails, err := _self.getFriendEmailsWithoutBlocking(ctx, secondUserID)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}

	// Get common friends
	commonFriends := make([]string, 0)
	commonMap := make(map[string]bool)
	for _, firstEmail := range firstFriendEmails {
		commonMap[firstEmail] = true
	}

	for _, secondEmail := range secondFriendEmails {
		if _, ok := commonMap[secondEmail]; ok {
			commonFriends = append(commonFriends, secondEmail)
		}
	}

	return commonFriends, nil
}

func (_self FriendService) CreateSubscription(ctx context.Context, requestorEmail string, targetEmail string) error {
	// Get requestor id and user target id from repository
	requestorId, err := _self.Repo.GetUserIDByEmail(ctx, requestorEmail)
	if err != nil {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: requestorEmail + " is not exists"}
	}
	targetId, err := _self.Repo.GetUserIDByEmail(ctx, targetEmail)
	if err != nil {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: targetEmail + " is not exists"}
	}

	// Check subscription relationship is exists
	isSubscribed, err := _self.Repo.IsSubscribedUser(ctx, requestorId, targetId)
	if err != nil {
		return &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}
	if isSubscribed {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: errs.MsgExistedSubscription}
	}

	// Check blocking between 2 user
	isBlocked, err := _self.Repo.IsBlockedUser(ctx, requestorId, targetId)
	if err != nil {
		return &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}
	if isBlocked {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: errs.MsgExistedBlockedUser}
	}

	if err := _self.Repo.CreateSubscription(ctx, requestorId, targetId); err != nil {
		return &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}

	return nil
}

func (_self FriendService) CreateUserBlock(ctx context.Context, requestorEmail string, targetEmail string) error {
	// Get requestor id and user target id from repository
	requestorId, err := _self.Repo.GetUserIDByEmail(ctx, requestorEmail)
	if err != nil {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: requestorEmail + " is not exists"}
	}
	targetId, err := _self.Repo.GetUserIDByEmail(ctx, targetEmail)
	if err != nil {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: targetEmail + " is not exists"}
	}

	// Check blocking between 2 user
	isBlocked, err := _self.Repo.IsBlockedUser(ctx, requestorId, targetId)
	if err != nil {
		return &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}
	if isBlocked {
		return &errs.FriendError{Code: http.StatusBadGateway, Description: errs.MsgExistedBlockedUser}
	}

	if err := _self.Repo.CreateUserBlock(ctx, requestorId, targetId); err != nil {
		return &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}

	return nil
}

func (_self FriendService) GetRecipientEmails(ctx context.Context, senderEmail string, text string) ([]string, error) {
	// Check existed email and get userID
	senderID, err := _self.Repo.GetUserIDByEmail(ctx, senderEmail)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusBadGateway, Description: senderEmail + " is not exists"}
	}

	recipients, err := _self.Repo.GetRecipientEmails(ctx, senderID)
	if err != nil {
		return nil, &errs.FriendError{Code: http.StatusInternalServerError, Description: err.Error()}
	}

	result := make([]string, 0)
	existedEmailsMap := make(map[string]bool)

	for _, user := range recipients {
		result = append(result, user.Email)
		existedEmailsMap[user.Email] = true
	}

	// Add mentioned emails to result
	mentionedEmails := utils.GetMentionedEmailFromText(text)
	for _, email := range mentionedEmails {
		if _, ok := existedEmailsMap[email]; !ok {
			result = append(result, email)
		}
	}

	return result, nil
}

// Get emails of users who are not being blocked by user
func (_self FriendService) getFriendEmailsWithoutBlocking(ctx context.Context, userId int) ([]string, error) {
	// Get friends by user id
	friendSlice, err := _self.Repo.GetFriendsByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	friendIds := make([]int, 0)
	for _, friend := range friendSlice {
		if friend.UserID == userId {
			friendIds = append(friendIds, friend.FriendID)
		}
		if friend.FriendID == userId {
			friendIds = append(friendIds, friend.UserID)
		}
	}

	// Get list users who have blocked user
	userBlocksSlice, err := _self.Repo.GetUserBlocksByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	userBlocksID := make([]int, 0)
	for _, user := range userBlocksSlice {
		if user.RequestorID == userId {
			userBlocksID = append(userBlocksID, user.TargetID)
		}
		if user.TargetID == userId {
			userBlocksID = append(userBlocksID, user.RequestorID)
		}
	}
	//Get UserID list with no blocked
	blockList := make(map[int]bool)
	for _, id := range userBlocksID {
		blockList[id] = true
	}
	friendIDsNonBlock := make([]int, 0)
	for _, id := range friendIds {
		if _, isBlock := blockList[id]; !isBlock {
			friendIDsNonBlock = append(friendIDsNonBlock, id)
		}
	}
	emails, err := _self.Repo.GetEmailsByUserIDs(ctx, friendIDsNonBlock)
	if err != nil {
		return nil, err
	}

	return emails, nil
}
