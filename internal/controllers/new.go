package controllers

import (
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/services"
)

type FriendController struct {
	Service services.SpecService
}

func NewFriendController(service services.SpecService) FriendController {
	return FriendController{
		Service: service,
	}
}
