package repository

import (
	"go-boilerplate/domain/user/interfaces"
	"go-boilerplate/model/entity"
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(ctx context.Context, id uint) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return user, err
}

func (r *repository) FindAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *repository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}
