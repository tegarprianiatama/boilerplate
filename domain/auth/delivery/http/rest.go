package http

import (
	"go-boilerplate/domain/auth/interfaces"
	"go-boilerplate/model/request"
	ce "go-boilerplate/shared/common/errors"
	jsonResp "go-boilerplate/shared/common/json"
	"errors"

	"github.com/labstack/echo/v4"
)

type authHandler struct {
	usecase interfaces.Usecase
	json    *jsonResp.Response
}

func AuthHandler(g *echo.Group, u interfaces.Usecase, j *jsonResp.Response) {
	h := &authHandler{
		usecase: u,
		json:    j,
	}

	authGroup := g.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
	}
}

// Register @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.CreateUserReq true "Register payload"
// @Success 201 {object} jsonResp.APIResponse
// @Failure 400 {object} jsonResp.APIResponse
// @Router /auth/register [post]
func (h *authHandler) Register(c echo.Context) error {
	var req request.CreateUserReq
	if err := c.Bind(&req); err != nil {
		return h.json.BadRequest(c, "Invalid request payload", err)
	}
	if err := c.Validate(&req); err != nil {
		return h.json.BadRequest(c, "Validation failed", err)
	}

	res, err := h.usecase.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, ce.ErrEmailAlreadyRegistered) {
			return h.json.BadRequest(c, err.Error(), err)
		}
		return h.json.InternalServerError(c, "Failed to register user", err)
	}

	return h.json.Created(c, "User registered successfully", res)
}

// Login @Summary Authenticate a user and issue a JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginReq true "Login credentials"
// @Success 200 {object} jsonResp.APIResponse
// @Failure 400 {object} jsonResp.APIResponse
// @Failure 401 {object} jsonResp.APIResponse
// @Router /auth/login [post]
func (h *authHandler) Login(c echo.Context) error {
	var req request.LoginReq
	if err := c.Bind(&req); err != nil {
		return h.json.BadRequest(c, "Invalid request payload", err)
	}
	if err := c.Validate(&req); err != nil {
		return h.json.BadRequest(c, "Validation failed", err)
	}

	res, err := h.usecase.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, ce.ErrInvalidCredentials) {
			return h.json.Unauthorized(c, err.Error())
		}
		return h.json.InternalServerError(c, "Failed to authenticate user", err)
	}

	return h.json.Ok(c, "Login successful", res)
}
