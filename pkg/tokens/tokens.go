package tokens

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

type Tokens struct {
	parser paseto.Parser
	secret paseto.V4AsymmetricSecretKey
	public paseto.V4AsymmetricPublicKey
}

var tokenKey string = "userid"

func (pst *Tokens) Verify(payload string) (string, error) {
	token, err := pst.parser.ParseV4Public(pst.public, payload, nil)
	if err != nil {
		return "", err
	}

	userID, err := token.GetString(tokenKey)
	return userID, err
}

func (pst *Tokens) CreateToken(payload string, exp time.Duration) string {
	token := paseto.NewToken()
	now := time.Now()
	token.SetIssuedAt(now)
	token.SetExpiration(now.Add(exp))
	token.SetString(tokenKey, payload)

	return token.V4Sign(pst.secret, nil)
}

func NewPaseto() *Tokens {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	secret := paseto.NewV4AsymmetricSecretKey()
	public := secret.Public()
	return &Tokens{parser, secret, public}
}
