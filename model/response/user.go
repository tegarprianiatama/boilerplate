package response

import (
	"boilerplate/model/entity"
	"time"
)

type UserResp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserResp(user entity.User) UserResp {
	return UserResp{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserRespList(users []entity.User) []UserResp {
	list := make([]UserResp, len(users))
	for i, u := range users {
		list[i] = ToUserResp(u)
	}
	return list
}
