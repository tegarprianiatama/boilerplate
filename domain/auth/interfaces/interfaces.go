package interfaces

import (
	"boilerplate/model/entity"
	"boilerplate/model/request"
	"boilerplate/model/response"
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
