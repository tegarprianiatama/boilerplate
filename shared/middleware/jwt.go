package middleware

import (
	jsonResp "go-boilerplate/shared/common/json"
	"go-boilerplate/shared/config"
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
				return jsonHelper.Unauthorized(c, "Missing or malformed JWT")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return jsonHelper.Unauthorized(c, "Invalid authorization format, expected Bearer {token}")
			}

			tokenStr := parts[1]
			claims := &jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || token == nil || !token.Valid {
				return jsonHelper.Unauthorized(c, "Invalid or expired JWT token")
			}

			// Store claims in context
			c.Set("user_claims", claims)
			return next(c)
		}
	}
}
