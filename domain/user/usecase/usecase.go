package usecase

import (
	"boilerplate/domain/user/interfaces"
	"boilerplate/model/entity"
	"boilerplate/model/request"
	"boilerplate/model/response"
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

type usecase struct {
	repo interfaces.Repository
}

func NewUsecase(repo interfaces.Repository) interfaces.Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(ctx context.Context, req request.CreateUserReq) (response.UserResp, error) {
	// Check email uniqueness
	existing, err := u.repo.FindByEmail(ctx, req.Email)
	if err == nil && existing.ID != 0 {
		return response.UserResp{}, errors.New("email already registered")
	}

	user := entity.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	if err := u.repo.Create(ctx, &user); err != nil {
		return response.UserResp{}, err
	}

	return response.ToUserResp(user), nil
}

func (u *usecase) GetByID(ctx context.Context, id uint) (response.UserResp, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.UserResp{}, errors.New("user not found")
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
			return response.UserResp{}, errors.New("user not found")
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
			return errors.New("user not found")
		}
		return err
	}
	return u.repo.Delete(ctx, id)
}
