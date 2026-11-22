package middleware

import (
	"boilerplate/shared/config"
	jsonResp "boilerplate/shared/common/json"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	jsonHelper := jsonResp.NewResponse()
	cfg := config.GetConfig()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return jsonHelper.BadRequest(c, "Missing or malformed JWT", "Authorization header is required")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return jsonHelper.BadRequest(c, "Invalid authorization format", "Format must be Bearer {token}")
			}

			tokenStr := parts[1]
			claims := &jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return jsonHelper.BadRequest(c, "Invalid or expired JWT token", err)
			}

			// Store claims in context
			c.Set("user_claims", claims)
			return next(c)
		}
	}
}
