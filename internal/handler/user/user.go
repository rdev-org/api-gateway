package user

import (
	pb "github.com/rdev-org/proto-contracts/pb/user"
)

type UserHandler struct {
	c pb.UserServiceClient
}

func NewUserService(c pb.UserServiceClient) *UserHandler {
	return &UserHandler{
		c: c,
	}
}
