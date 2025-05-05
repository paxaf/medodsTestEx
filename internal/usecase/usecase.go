package usecase

import (
	"github.com/paxaf/medodsTestEx/internal/tokens"
)

func GetTokens(guid string) (*tokens.Tokens, error) {
	tokenAll := tokens.NewTokens(guid)
	hashedRefresh := tokenAll.GetHashedRefresh()
}
