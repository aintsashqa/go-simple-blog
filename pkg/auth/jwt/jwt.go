package jwt

import (
	"errors"
	"time"

	"github.com/aintsashqa/go-simple-blog/pkg/auth"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type JWTAuthorizationProvider struct {
	signingKey string
}

func NewJWTAuthorizationProvider(signingKey string) *JWTAuthorizationProvider {
	return &JWTAuthorizationProvider{signingKey: signingKey}
}

func (p *JWTAuthorizationProvider) NewToken(params auth.TokenParams) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   params.UserID.String(),
		ExpiresAt: time.Now().Add(params.ExpiresAt).Unix(),
	})
	return token.SignedString([]byte(p.signingKey))
}

func (p *JWTAuthorizationProvider) Parse(value string) (uuid.UUID, error) {
	token, err := jwt.Parse(value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(p.signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("error get user claims from token")
	}

	return uuid.FromStringOrNil(claims["sub"].(string)), nil
}
