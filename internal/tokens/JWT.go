package tokens

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret")

type Tokens struct {
	accessToken  string
	refreshToken string
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
	refresh, err := GenerateRefresh(guid)
	if err != nil {
		return nil
	}
	return &Tokens{
		accessToken:  access,
		refreshToken: refresh,
	}
}

func GenerateJWT(guid string) (string, error) {
	claims := Claims{
		UserGUID: guid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}

func (t *Tokens) GetHashedRefresh() string {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(t.refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err)
		return ""
	}
	return string(hashedToken)
}
