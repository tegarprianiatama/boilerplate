package usecase

import (
	"boilerplate/domain/user/interfaces"
	"boilerplate/model/entity"
	"boilerplate/model/request"
	"boilerplate/model/response"
	"boilerplate/shared/auth"
	ce "boilerplate/shared/common/errors"
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

func (u *usecase) Create(ctx context.Context, req request.CreateUserReq) (response.UserResp, error) {
	existing, err := u.repo.FindByEmail(ctx, req.Email)
	if err == nil && existing.ID != 0 {
		return response.UserResp{}, ce.ErrEmailAlreadyRegistered
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response.UserResp{}, err
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return response.UserResp{}, err
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
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

func (u *usecase) Login(ctx context.Context, req request.LoginReq) (response.LoginResp, error) {
	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.LoginResp{}, ce.ErrInvalidCredentials
		}
		return response.LoginResp{}, err
	}

	if !auth.CheckPassword(user.Password, req.Password) {
		return response.LoginResp{}, ce.ErrInvalidCredentials
	}

	token, expiresAt, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return response.LoginResp{}, err
	}

	return response.LoginResp{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
		User:      response.ToUserResp(user),
	}, nil
}
