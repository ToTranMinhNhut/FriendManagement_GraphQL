package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/graph/generated"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/graph/graphmodel"
)

func (r *mutationResolver) CreateFriend(ctx context.Context, input graphmodel.Friends) (*graphmodel.IsSuccess, error) {
	// Decode request body
	friendReq := FriendRequest{}
	for _, email := range input.Friends {
		friendReq.Emails = append(friendReq.Emails, email)
	}

	//Validation
	if err := friendReq.Validate(); err != nil {
		return nil, err
	}

	if err := r.Service.CreateFriend(ctx, friendReq.Emails[0], friendReq.Emails[1]); err != nil {
		return nil, err
	}

	//Response
	return &graphmodel.IsSuccess{
		Success: true,
	}, nil
}

func (r *mutationResolver) FriendList(ctx context.Context, input graphmodel.Email) (*graphmodel.FriendList, error) {
	//Decode request body
	userReq := UserRequest{
		Email: input.Email,
	}

	//Validation
	if err := userReq.Validate(); err != nil {
		return nil, err
	}

	friendEmails, err := r.Service.GetFriends(ctx, userReq.Email)
	if err != nil {
		return nil, err
	}

	//Response
	return &graphmodel.FriendList{
		Success: true,
		Friends: friendEmails,
		Count:   len(friendEmails),
	}, nil
}

func (r *mutationResolver) CommonFriends(ctx context.Context, input graphmodel.Friends) (*graphmodel.FriendList, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Subscribe(ctx context.Context, input graphmodel.RequestTarget) (*graphmodel.IsSuccess, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) BlockUpdate(ctx context.Context, input graphmodel.RequestTarget) (*graphmodel.IsSuccess, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RetrieveEmailReceiveUpdate(ctx context.Context, input graphmodel.SendMail) (*graphmodel.Recipients, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) (*graphmodel.Users, error) {
	emails, err := r.Service.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	//Response
	return &graphmodel.Users{
		Success: true,
		Emails:  emails,
		Count:   len(emails),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
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
