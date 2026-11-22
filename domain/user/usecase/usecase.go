package usecase

import (
	"boilerplate/domain/user/interfaces"
	ce "boilerplate/shared/common/errors"
	"boilerplate/model/request"
	"boilerplate/model/response"
	"context"
	"errors"

	"gorm.io/gorm"
)

type usecase struct {
	repo interfaces.Repository
}

func NewUsecase(repo interfaces.Repository) interfaces.Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) GetByID(ctx context.Context, id uint) (response.UserResp, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.UserResp{}, ce.ErrUserNotFound
		}
		return response.UserResp{}, err
	}
	return response.ToUserResp(user), nil
}

func (u *usecase) GetAll(ctx context.Context) ([]response.UserResp, error) {
	users, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return response.ToUserRespList(users), nil
}

func (u *usecase) Update(ctx context.Context, id uint, req request.UpdateUserReq) (response.UserResp, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.UserResp{}, ce.ErrUserNotFound
		}
		return response.UserResp{}, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := u.repo.Update(ctx, &user); err != nil {
		return response.UserResp{}, err
	}

	return response.ToUserResp(user), nil
}

func (u *usecase) Delete(ctx context.Context, id uint) error {
	_, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ce.ErrUserNotFound
		}
		return err
	}
	return u.repo.Delete(ctx, id)
}
