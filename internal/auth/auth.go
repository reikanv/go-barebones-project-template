package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/bobopylabepolhk/ypgophermart/internal/db"
	"github.com/bobopylabepolhk/ypgophermart/pkg/tokens"
	"github.com/jackc/pgx/v5"
)

type (
	auth struct {
		db     *db.Db
		tokens *tokens.Tokens
	}
)

func generateSalt() (string, []byte, error) {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		return "", nil, err
	}

	return string(bs), bs, nil
}

func hashPassword(password string, bsalt []byte) string {
	// SHA256 hash the password and append salt
	bpassword := []byte(password)
	hasher := sha256.New()
	hasher.Write(bpassword)
	bhash := hasher.Sum(bsalt)

	return string(bhash)
}

func (s *auth) Login(ctx context.Context, login, password string) (int, error) {
	row := s.db.Pool.QueryRow(ctx, "select id, password, salt from user where login = $1", login)
	var id int
	var storedPassword string
	var salt string
	if err := row.Scan(&id, &storedPassword, &salt); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return 0, errWrongCredentials
		}
		return 0, fmt.Errorf("auth - Login: %w", err)
	}

	bsalt := []byte(salt)
	hashedPassword := hashPassword(password, bsalt)
	if hashedPassword != storedPassword {
		return 0, errWrongCredentials
	}

	return id, nil
}

func (s *auth) CreateUser(ctx context.Context, login, password string) (int, error) {
	salt, bsalt, err := generateSalt()
	if err != nil {
		return 0, fmt.Errorf("auth - CreateUser crypto failed %w", err)
	}

	row := s.db.Pool.QueryRow(
		ctx,
		"insert into user (login, password, salt) VALUES ($1, $2, $3) on conflict do nothing returning id",
		login,
		hashPassword(password, bsalt),
		salt,
	)

	var id int
	if err := row.Scan(&id); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return 0, errUserAlreadyExists(login)
		}
		return 0, fmt.Errorf("auth - CreateUser: %w", err)
	}
	return id, nil
}

func (s *auth) GrantAccessToken(id int, expiresIn time.Duration) string {
	return s.tokens.CreateToken(fmt.Sprint(id), expiresIn)
}

func NewAuthService(db *db.Db, t *tokens.Tokens) *auth {
	return &auth{db, t}
}
