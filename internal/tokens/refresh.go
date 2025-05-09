package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRefresh() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	base64Token := base64.StdEncoding.EncodeToString([]byte(b))

	return base64Token, nil
}

func (t *Tokens) ValidateRefresh(hashedToken string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(t.RefreshToken))
	fmt.Println("CHAP err: ", err)
	return err == nil
}

func (t *Tokens) GetHashedRefresh() (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(t.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	return string(hashedToken), nil
}
