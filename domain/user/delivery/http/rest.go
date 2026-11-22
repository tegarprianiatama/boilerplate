package http

import (
	"boilerplate/domain/user/interfaces"
	"boilerplate/model/request"
	ce "boilerplate/shared/common/errors"
	jsonResp "boilerplate/shared/common/json"
	"boilerplate/shared/middleware"
	"errors"
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

	userGroup := g.Group("/users", middleware.JWTMiddleware())
	{
		userGroup.GET("", h.GetAll)
		userGroup.GET("/:id", h.GetByID)
		userGroup.PUT("/:id", h.Update)
		userGroup.DELETE("/:id", h.Delete)
	}
}

// GetAll @Summary Get all users
// @Tags users
// @Produce json
// @Success 200 {object} jsonResp.APIResponse
// @Failure 500 {object} jsonResp.APIResponse
// @Router /users [get]
func (h *userHandler) GetAll(c echo.Context) error {
	res, err := h.usecase.GetAll(c.Request().Context())
	if err != nil {
		return h.json.InternalServerError(c, "Failed to retrieve users", err)
	}
	return h.json.Ok(c, "Users retrieved successfully", res)
}

// GetByID @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} jsonResp.APIResponse
// @Failure 400 {object} jsonResp.APIResponse
// @Failure 404 {object} jsonResp.APIResponse
// @Router /users/{id} [get]
func (h *userHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return h.json.BadRequest(c, "Invalid user ID", err)
	}

	res, err := h.usecase.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		if errors.Is(err, ce.ErrUserNotFound) {
			return h.json.NotFound(c, err.Error())
		}
		return h.json.InternalServerError(c, "Failed to retrieve user", err)
	}

	return h.json.Ok(c, "User retrieved successfully", res)
}

// Update @Summary Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body request.UpdateUserReq true "Update user payload"
// @Success 200 {object} jsonResp.APIResponse
// @Failure 400 {object} jsonResp.APIResponse
// @Failure 404 {object} jsonResp.APIResponse
// @Router /users/{id} [put]
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
		if errors.Is(err, ce.ErrUserNotFound) {
			return h.json.NotFound(c, err.Error())
		}
		return h.json.BadRequest(c, err.Error(), err)
	}

	return h.json.Ok(c, "User updated successfully", res)
}

// Delete @Summary Delete user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} jsonResp.APIResponse
// @Failure 400 {object} jsonResp.APIResponse
// @Failure 404 {object} jsonResp.APIResponse
// @Router /users/{id} [delete]
func (h *userHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return h.json.BadRequest(c, "Invalid user ID", err)
	}

	err = h.usecase.Delete(c.Request().Context(), uint(id))
	if err != nil {
		if errors.Is(err, ce.ErrUserNotFound) {
			return h.json.NotFound(c, err.Error())
		}
		return h.json.InternalServerError(c, "Failed to delete user", err)
	}

	return h.json.Ok(c, "User deleted successfully", nil)
}
