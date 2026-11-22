package interfaces

import (
	"go-boilerplate/model/entity"
	"go-boilerplate/model/request"
	"go-boilerplate/model/response"
	"context"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

type Usecase interface {
	Register(ctx context.Context, req request.CreateUserReq) (response.UserResp, error)
	Login(ctx context.Context, req request.LoginReq) (response.LoginResp, error)
}
