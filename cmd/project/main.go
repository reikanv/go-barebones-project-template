package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"

	"github.com/bobopylabepolhk/ypgophermart/config"
	"github.com/bobopylabepolhk/ypgophermart/internal/auth"
	"github.com/bobopylabepolhk/ypgophermart/internal/db"
	"github.com/bobopylabepolhk/ypgophermart/pkg/logger"
	"github.com/bobopylabepolhk/ypgophermart/pkg/middleware"
	"github.com/bobopylabepolhk/ypgophermart/pkg/tokens"
)

func setupRouters(e *echo.Echo, db *db.Db, t *tokens.Tokens, cfg *config.Config) {
	baseGroup := e.Group(cfg.RootPath)

	auth.NewAuthRouter(baseGroup, db, t, cfg.AuthCookieName)
}

func setupMiddleware(e *echo.Echo, t *tokens.Tokens, cfg *config.Config) {
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())
	e.Use(middleware.CheckAuth(t, cfg.AuthCookieName))
}

func run(cfg *config.Config) {
	pg, err := db.NewPGX(context.Background(), cfg.DatabaseURI)
	if err != nil {
		slog.Error(err.Error())
	}
	defer pg.Close()
	slog.Info("connected to db")

	t := tokens.NewPaseto()

	e := echo.New()
	e.HideBanner = true
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	setupMiddleware(e, t, cfg)
	setupRouters(e, pg, t, cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	go func() {
		err := e.Start(cfg.RunAddress)
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	// shutdown on interrupt
	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		slog.Error(err.Error())
	}

	logger.InitLogger(cfg.Debug)

	run(cfg)
}
