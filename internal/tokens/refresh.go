package tokens

import (
	"encoding/base64"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRefresh(guid string) (string, error) {
	base64Token := base64.StdEncoding.EncodeToString([]byte(guid))

	return base64Token, nil
}

func ValidateRefresh(clientToken, hashedToken string) bool {
	decodeToken, err := base64.StdEncoding.DecodeString(clientToken)
	if err != nil {
		log.Info().Msg("invalid token")
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedToken), decodeToken)
	return err == nil
}
