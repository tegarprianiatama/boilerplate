package repository

import (
	"boilerplate/domain/user/interfaces"
	"boilerplate/model/entity"
	"context"

	"github.com/jinzhu/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *repository) FindAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) Update(ctx context.Context, user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}
