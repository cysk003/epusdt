package middleware

import (
	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/util/constant"
	"github.com/labstack/echo/v4"
)

// CheckApiToken validates the Authorization header against the configured api_auth_token.
// Use this for management APIs (GET/POST) where body-based signature is not practical.
func CheckApiToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			token := ctx.Request().Header.Get("Authorization")
			if token == "" {
				token = ctx.QueryParam("api_token")
			}
			if token == "" || token != config.GetApiAuthToken() {
				return constant.SignatureErr
			}
			return next(ctx)
		}
	}
}
