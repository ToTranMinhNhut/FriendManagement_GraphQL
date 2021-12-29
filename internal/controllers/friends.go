package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/errs"
)

type FriendRequest struct {
	Emails []string `json:"friends"`
}

type UserRequest struct {
	Email string `json:"email"`
}

type RequestorRequest struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

type RecipientsRequest struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

// Get all of users
func (_self FriendController) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// Validate body request
	if r.ContentLength != 0 {
		Respond(w, http.StatusBadRequest, MsgError(errs.ErrBodyRequestInvalid))
		return
	}

	emails, err := _self.Service.GetUsers(ctx)
	if err != nil {
		if friendErr, ok := err.(*errs.FriendError); ok && friendErr != nil {
			Respond(w, friendErr.Code, MsgError(friendErr))
			return
		}
		Respond(w, http.StatusInternalServerError, MsgError(err))
		return
	}

	Respond(w, http.StatusOK, MsgGetAllUsersOk(emails, len(emails)))
}

// Create a new friend relationship
func (_self FriendController) CreateFriend(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	friendReq := FriendRequest{}
	if err := json.NewDecoder(r.Body).Decode(&friendReq); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(errs.ErrBodyRequestInvalid))
		return
	}

	// Validate request body
	if err := friendReq.Validate(); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(err))
		return
	}

	if err := _self.Service.CreateFriend(ctx, friendReq.Emails[0], friendReq.Emails[1]); err != nil {
		if friendErr, ok := err.(*errs.FriendError); ok && friendErr != nil {
			Respond(w, friendErr.Code, MsgError(friendErr))
			return
		}
		Respond(w, http.StatusInternalServerError, MsgError(err))
		return
	}

	Respond(w, http.StatusOK, MsgOK())
}

// Get all of friends of a user without blocking relationship
func (_self FriendController) GetFriends(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userReq := UserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(errs.ErrBodyRequestInvalid))
		return
	}

	// Validation request body
	if err := userReq.Validate(); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(err))
		return
	}

	friendEmails, err := _self.Service.GetFriends(ctx, userReq.Email)
	if err != nil {
		if friendErr, ok := err.(*errs.FriendError); ok && friendErr != nil {
			Respond(w, friendErr.Code, MsgError(friendErr))
			return
		}
		Respond(w, http.StatusInternalServerError, MsgError(err))
		return
	}

	Respond(w, http.StatusOK, MsgGetFriendsOk(friendEmails, len(friendEmails)))
}

// Get common friends of 2 users
func (_self FriendController) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	friendReq := FriendRequest{}
	if err := json.NewDecoder(r.Body).Decode(&friendReq); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(errs.ErrBodyRequestInvalid))
		return
	}

	// Validate request body
	if err := friendReq.Validate(); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(err))
		return
	}

	// Get common friends
	commonFriends, err := _self.Service.GetCommonFriends(ctx, friendReq.Emails[0], friendReq.Emails[1])
	if err != nil {
		if friendErr, ok := err.(*errs.FriendError); ok && friendErr != nil {
			Respond(w, friendErr.Code, MsgError(friendErr))
			return
		}
		Respond(w, http.StatusInternalServerError, MsgError(err))
		return
	}

	Respond(w, http.StatusOK, MsgGetFriendsOk(commonFriends, len(commonFriends)))
}

// Create a subscription relationship of users
func (_self FriendController) CreateSubcription(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	requestorReq := RequestorRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestorReq); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(errs.ErrBodyRequestInvalid))
		return
	}

	// Validate request
	if err := requestorReq.Validate(); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(err))
		return
	}

	if err := _self.Service.CreateSubscription(ctx, requestorReq.Requestor, requestorReq.Target); err != nil {
		if friendErr, ok := err.(*errs.FriendError); ok && friendErr != nil {
			Respond(w, friendErr.Code, MsgError(friendErr))
			return
		}
		Respond(w, http.StatusInternalServerError, MsgError(err))
		return
	}

	Respond(w, http.StatusOK, MsgOK())
}

// Create a blocking relationship of users
func (_self FriendController) CreateUserBlock(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	requestorReq := RequestorRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestorReq); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(errs.ErrBodyRequestInvalid))
		return
	}

	// Validate request
	if err := requestorReq.Validate(); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(err))
		return
	}

	if err := _self.Service.CreateUserBlock(ctx, requestorReq.Requestor, requestorReq.Target); err != nil {
		if friendErr, ok := err.(*errs.FriendError); ok && friendErr != nil {
			Respond(w, friendErr.Code, MsgError(friendErr))
			return
		}
		Respond(w, http.StatusInternalServerError, MsgError(err))
		return
	}

	Respond(w, http.StatusOK, MsgOK())
}

// Get all of recipients who are friend, subscriber, and mention user without blocking by user
func (_self FriendController) GetRecipientEmails(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	recipient := RecipientsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&recipient); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(err))
		return
	}

	// Validate request body
	if err := recipient.Validate(); err != nil {
		Respond(w, http.StatusBadRequest, MsgError(err))
		return
	}

	recipients, err := _self.Service.GetRecipientEmails(ctx, recipient.Sender, recipient.Text)
	if err != nil {
		if friendErr, ok := err.(*errs.FriendError); ok && friendErr != nil {
			Respond(w, friendErr.Code, MsgError(friendErr))
			return
		}
		Respond(w, http.StatusInternalServerError, MsgError(err))
		return
	}

	Respond(w, http.StatusOK, MsgGetEmailReceiversOk(recipients))
}
