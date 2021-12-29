package services

import "github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/repository"

type FriendService struct {
	Repo repository.SpecRepo
}

func NewFriendService(repo repository.SpecRepo) FriendService {
	return FriendService{
		Repo: repo,
	}
}
