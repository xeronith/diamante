package jwt

import (
	"errors"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt"
)

type TokenPayload struct {
	ID int64
}

func Generate(tokenKey, tokenExpiration string) string {
	payload := &TokenPayload{ID: time.Now().UnixNano()}

	v, err := time.ParseDuration(tokenExpiration)

	if err != nil {
		panic("Invalid time duration. Should be time.ParseDuration string")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(v).Unix(),
		"ID":  payload.ID,
	})

	token, err := t.SignedString([]byte(tokenKey))

	if err != nil {
		panic(err)
	}

	return token
}

func parse(token, tokenKey string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(tokenKey), nil
	})
}

func Verify(token, tokenKey string) (*TokenPayload, error) {
	parsed, err := parse(token, tokenKey)

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
