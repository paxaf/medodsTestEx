package usecase

import (
	"github.com/paxaf/medodsTestEx/internal/repository"
	"github.com/paxaf/medodsTestEx/internal/tokens"
)

type UseCase struct {
	repo *repository.Repository
}

func NewUseCase(repo *repository.Repository) UseCase {
	return UseCase{
		repo: repo,
	}
}
func (u *UseCase) GetTokens(guid string) (*tokens.Tokens, error) {
	tokenAll := tokens.NewTokens(guid)
	hashedRefresh, err := tokenAll.GetHashedRefresh()
	if err != nil {
		return nil, err
	}
	err = u.repo.SetHash(guid, hashedRefresh)
	if err != nil {
		return nil, err
	}
	return tokenAll, nil
}
