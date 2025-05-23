package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/bobopylabepolhk/ypgophermart/pkg/tokens"
)

func CheckAuth(t *tokens.Tokens, cookieName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cookie, err := ctx.Cookie(cookieName)
			if err != nil {
				return echo.ErrUnauthorized
			}

			if userID, err := t.Verify(cookie.Value); err == nil {
				ctx.Set(cookieName, userID)
				return next(ctx)
			}

			return echo.ErrUnauthorized
		}
	}
}
