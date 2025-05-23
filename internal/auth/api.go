package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/bobopylabepolhk/ypgophermart/internal/db"
	"github.com/bobopylabepolhk/ypgophermart/pkg/tokens"
	"github.com/labstack/echo/v4"
)

type (
	authRouter struct {
		AuthService AuthService
		CookieName  string
	}
	AuthService interface {
		CreateUser(ctx context.Context, login, password string) (int, error)
		Login(ctx context.Context, login, password string) (int, error)
		GrantAccessToken(id int, expiresIn time.Duration) string
	}
	credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
)

var expiresIn time.Duration = time.Hour * 24

func createCookie(accessToken string, cookieName string) *http.Cookie {
	c := new(http.Cookie)
	c.Name = cookieName
	c.Value = accessToken
	c.Expires = time.Now().Add(expiresIn)
	return c
}

func (r *authRouter) loginHandler(ctx echo.Context) error {
	var cred credentials
	if err := ctx.Bind(&cred); err != nil {
		slog.Error(err.Error())
		return echo.ErrBadRequest
	}

	id, err := r.AuthService.Login(ctx.Request().Context(), cred.Login, cred.Password)
	if errors.Is(err, errWrongCredentials) {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	token := r.AuthService.GrantAccessToken(id, expiresIn)
	ctx.SetCookie(createCookie(token, r.CookieName))

	return ctx.NoContent(http.StatusOK)
}

func (r *authRouter) registerHandler(ctx echo.Context) error {
	var cred credentials
	if err := ctx.Bind(&cred); err != nil {
		slog.Error(err.Error())
		return echo.ErrBadRequest
	}

	id, err := r.AuthService.CreateUser(ctx.Request().Context(), cred.Login, cred.Password)
	if errors.Is(err, errUserAlreadyExists(cred.Login)) {
		return ctx.NoContent(http.StatusConflict)
	}

	token := r.AuthService.GrantAccessToken(id, expiresIn)
	ctx.SetCookie(createCookie(token, r.CookieName))

	return ctx.NoContent(http.StatusOK)
}

func NewAuthRouter(baseGroup *echo.Group, db *db.Db, t *tokens.Tokens, cookieName string) {
	authService := NewAuthService(db, t)
	router := &authRouter{AuthService: authService, CookieName: cookieName}

	baseGroup.POST("/login", router.loginHandler)
	baseGroup.POST("/register", router.registerHandler)
}
