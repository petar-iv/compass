package jwtutil

import (
	"errors"

	"github.com/form3tech-oss/jwt-go"
)

type Claims struct {
	ClientId string `json:"client_id"`
	Scopes   string `json:"scope"`
}

func NewClaims(bearerToken string) (Claims, error) {
	claims := Claims{}

	s := jwt.Parser{SkipClaimsValidation: true}
	s.ParseUnverified(bearerToken, &claims)

	return claims, nil
}

func (c Claims) Valid() error {
	if c.ClientId == "" {
		return errors.New("client id cannot be empty")
	}

	return nil
}
