package interfaces

import (
	"go-boilerplate/model/entity"
	"go-boilerplate/model/request"
	"go-boilerplate/model/response"
	"context"
)

type Repository interface {
	FindByID(ctx context.Context, id uint) (entity.User, error)
	FindAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
}

type Usecase interface {
	GetByID(ctx context.Context, id uint) (response.UserResp, error)
	GetAll(ctx context.Context) ([]response.UserResp, error)
	Update(ctx context.Context, id uint, req request.UpdateUserReq) (response.UserResp, error)
	Delete(ctx context.Context, id uint) error
}
