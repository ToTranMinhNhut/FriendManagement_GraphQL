package errs

import (
	"errors"
	"fmt"
)

var (
	ErrBodyRequestInvalid    = errors.New("Body request invalid format")
	ErrBodyRequestEmpty      = errors.New("Request body is empty")
	ErrNumberOfEmail         = errors.New("Number of email addresses must be 2")
	ErrDifferentEmail        = errors.New("Two email addresses must be different")
	ErrRequestorFieldInvalid = errors.New("Requestor field invalid format")
	ErrTargetFieldInvalid    = errors.New("Target field invalid format")
	ErrSenderFieldInvalid    = errors.New("Sender field invalid format")
	ErrTextFieldInvalid      = errors.New("Text field invalid format")

	MsgExistedFriendship   = "The friend relationship has been existed"
	MsgExistedBlockedUser  = "The users have blocked each other"
	MsgExistedSubscription = "The users have subscribed each other"
	MsgCreatedFriendship   = "Users cannot be created a new friendship"
)

type FriendError struct {
	Code        int    `json:"-"`
	Description string `json:"error_description"`
}

func (e *FriendError) Error() string {
	return fmt.Sprintf("%s", e.Description)
}
