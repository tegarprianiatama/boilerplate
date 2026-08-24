package http

import (
	"boilerplate/domain/user/interfaces"
	"boilerplate/model/request"
	jsonResp "boilerplate/shared/common/json"
	"strconv"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	usecase interfaces.Usecase
	json    *jsonResp.Response
}

func UserHandler(g *echo.Group, u interfaces.Usecase, j *jsonResp.Response) {
	h := &userHandler{
		usecase: u,
		json:    j,
	}

	userGroup := g.Group("/users")
	{
		userGroup.POST("", h.Create)
		userGroup.GET("", h.GetAll)
		userGroup.GET("/:id", h.GetByID)
		userGroup.PUT("/:id", h.Update)
		userGroup.DELETE("/:id", h.Delete)
	}
}

func (h *userHandler) Create(c echo.Context) error {
	var req request.CreateUserReq
	if err := c.Bind(&req); err != nil {
		return h.json.BadRequest(c, "Invalid request payload", err)
	}
	if err := c.Validate(&req); err != nil {
		return h.json.BadRequest(c, "Validation failed", err)
	}

	res, err := h.usecase.Create(c.Request().Context(), req)
	if err != nil {
		return h.json.BadRequest(c, err.Error(), err)
	}

	return h.json.Created(c, "User created successfully", res)
}

func (h *userHandler) GetAll(c echo.Context) error {
	res, err := h.usecase.GetAll(c.Request().Context())
	if err != nil {
		return h.json.InternalServerError(c, "Failed to retrieve users", err)
	}
	return h.json.Ok(c, "Users retrieved successfully", res)
}

func (h *userHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return h.json.BadRequest(c, "Invalid user ID", err)
	}

	res, err := h.usecase.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			return h.json.NotFound(c, err.Error())
		}
		return h.json.InternalServerError(c, "Failed to retrieve user", err)
	}

	return h.json.Ok(c, "User retrieved successfully", res)
}

func (h *userHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return h.json.BadRequest(c, "Invalid user ID", err)
	}

	var req request.UpdateUserReq
	if err := c.Bind(&req); err != nil {
		return h.json.BadRequest(c, "Invalid request payload", err)
	}
	if err := c.Validate(&req); err != nil {
		return h.json.BadRequest(c, "Validation failed", err)
	}

	res, err := h.usecase.Update(c.Request().Context(), uint(id), req)
	if err != nil {
		if err.Error() == "user not found" {
			return h.json.NotFound(c, err.Error())
		}
		return h.json.BadRequest(c, err.Error(), err)
	}

	return h.json.Ok(c, "User updated successfully", res)
}

func (h *userHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return h.json.BadRequest(c, "Invalid user ID", err)
	}

	err = h.usecase.Delete(c.Request().Context(), uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			return h.json.NotFound(c, err.Error())
		}
		return h.json.InternalServerError(c, "Failed to delete user", err)
	}

	return h.json.Ok(c, "User deleted successfully", nil)
}
