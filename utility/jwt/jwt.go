package jwt

import (
	"errors"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt"
)

const (
	// TOKEN_KEY returns the jwt token secret
	TOKEN_KEY = "A9VxLvehJEnwMcsfG18=atlXh$FMOe^7M619oDs1g=bcrTWgGNYpnowkog5rA"
	// TOKEN_EXP returns the jwt token expiration duration.
	// Should be time.ParseDuration string. Source: https://golang.org/pkg/time/#ParseDuration
	// default: 10h
	TOKEN_EXP = "10h"
)

type TokenPayload struct {
	ID int64
}

func Generate() string {
	payload := &TokenPayload{ID: time.Now().UnixNano()}

	v, err := time.ParseDuration(TOKEN_EXP)

	if err != nil {
		panic("Invalid time duration. Should be time.ParseDuration string")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(v).Unix(),
		"ID":  payload.ID,
	})

	token, err := t.SignedString([]byte(TOKEN_KEY))

	if err != nil {
		panic(err)
	}

	return token
}

func parse(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(TOKEN_KEY), nil
	})
}

func Verify(token string) (*TokenPayload, error) {
	parsed, err := parse(token)

	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	id, ok := claims["ID"].(float64)
	if !ok {
		return nil, errors.New("something went wrong")
	}

	return &TokenPayload{
		ID: int64(id),
	}, nil
}
