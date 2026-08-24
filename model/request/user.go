package request

type CreateUserReq struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin user manager"`
}

type UpdateUserReq struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
	Role string `json:"role" validate:"omitempty,oneof=admin user manager"`
}
