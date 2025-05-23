package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Gzip() echo.MiddlewareFunc {
	cfg := middleware.GzipConfig{
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			if strings.Contains(req.URL.Path, "swagger") {
				return true
			}

			return false
		},
	}

	return middleware.GzipWithConfig(cfg)
}

func Decompress() echo.MiddlewareFunc {
	cfg := middleware.DecompressConfig{
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			if strings.Contains(req.URL.Path, "swagger") {
				return true
			}

			return false
		},
	}
	return middleware.DecompressWithConfig(cfg)
}
