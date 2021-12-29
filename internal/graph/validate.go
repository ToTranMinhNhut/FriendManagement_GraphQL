package graph

import (
	"errors"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/errs"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/utils"
)

// Validate to body of friend request
func (_self FriendRequest) Validate() error {
	if _self.Emails == nil && len(_self.Emails) == 0 {
		return errs.ErrBodyRequestEmpty
	}
	if len(_self.Emails) != 2 {
		return errs.ErrNumberOfEmail
	}
	if _self.Emails[0] == _self.Emails[1] {
		return errs.ErrDifferentEmail
	}
	isValidUserEmail, err := utils.IsValidEmail(_self.Emails[0])
	if !isValidUserEmail || err != nil {
		return errors.New(_self.Emails[0] + " invalid format (ex: \"andy@example.com\")")
	}
	isValidFriendEmail, err := utils.IsValidEmail(_self.Emails[1])
	if !isValidFriendEmail || err != nil {
		return errors.New(_self.Emails[1] + " invalid format (ex: \"andy@example.com\")")
	}
	return nil
}

// Validate to body of user request
func (_self UserRequest) Validate() error {
	if _self.Email == "" {
		return errs.ErrBodyRequestEmpty
	}
	isValidEmail, err := utils.IsValidEmail(_self.Email)
	if !isValidEmail || err != nil {
		return errors.New(_self.Email + " invalid format (ex: \"andy@example.com\")")
	}
	return nil
}

// Validate to body of requestor request
func (_self RequestorRequest) Validate() error {
	if _self.Requestor == "" && _self.Target == "" {
		return errs.ErrBodyRequestEmpty
	}

	if _self.Requestor == "" {
		return errs.ErrRequestorFieldInvalid
	}

	if _self.Target == "" {
		return errs.ErrTargetFieldInvalid
	}

	if _self.Target == _self.Requestor {
		return errs.ErrDifferentEmail
	}

	isValidRequestEmail, requestErr := utils.IsValidEmail(_self.Requestor)
	if !isValidRequestEmail || requestErr != nil {
		return errors.New(_self.Requestor + " invalid format (ex: \"andy@example.com\")")
	}

	isValidTargetEmail, targetErr := utils.IsValidEmail(_self.Target)
	if !isValidTargetEmail || targetErr != nil {
		return errors.New(_self.Target + " invalid format (ex: \"andy@example.com\")")
	}
	return nil
}

// Validate to body of recipient request
func (_self RecipientsRequest) Validate() error {
	if _self.Sender == "" && _self.Text == "" {
		return errs.ErrBodyRequestEmpty
	}

	if _self.Sender == "" {
		return errs.ErrSenderFieldInvalid
	}
	if _self.Text == "" {
		return errs.ErrTextFieldInvalid
	}
	isValidEmail, err := utils.IsValidEmail(_self.Sender)

	if !isValidEmail || err != nil {
		return errors.New(_self.Sender + " invalid format (ex: \"andy@example.com\")")
	}
	return nil
}
