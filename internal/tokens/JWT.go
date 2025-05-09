package tokens

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret")

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Claims struct {
	UserGUID string `json:"user_guid"`
	jwt.RegisteredClaims
}

func NewTokens(guid string) *Tokens {
	access, err := GenerateJWT(guid)
	if err != nil {
		return nil
	}
	refresh, err := GenerateRefresh()
	if err != nil {
		return nil
	}
	return &Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func GenerateJWT(guid string) (string, error) {
	claims := Claims{
		UserGUID: guid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}

func (t *Tokens) ValidateJWT() (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(t.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.UserGUID, nil
}
