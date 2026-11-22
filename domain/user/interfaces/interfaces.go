package interfaces

import (
	"boilerplate/model/entity"
	"boilerplate/model/request"
	"boilerplate/model/response"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uint) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
}

type Usecase interface {
	Create(ctx context.Context, req request.CreateUserReq) (response.UserResp, error)
	GetByID(ctx context.Context, id uint) (response.UserResp, error)
	GetAll(ctx context.Context) ([]response.UserResp, error)
	Update(ctx context.Context, id uint, req request.UpdateUserReq) (response.UserResp, error)
	Delete(ctx context.Context, id uint) error
	Login(ctx context.Context, req request.LoginReq) (response.LoginResp, error)
}
