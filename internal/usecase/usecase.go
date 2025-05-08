package usecase

import (
	"github.com/google/uuid"
	"github.com/paxaf/medodsTestEx/internal/repository"
	"github.com/paxaf/medodsTestEx/internal/tokens"
)

type UseCase interface {
	GetTokens(guid string) (*tokens.Tokens, error)
}

type usecase struct {
	repo repository.PgRepository
}

func NewUseCase(repo repository.PgRepository) UseCase {
	return &usecase{
		repo: repo,
	}
}
func (u *usecase) GetTokens(guid string) (*tokens.Tokens, error) {
	err := uuid.Validate(guid)
	if err != nil {
		return nil, err
	}
	tokenAll := tokens.NewTokens(guid)
	hashedRefresh, err := tokenAll.GetHashedRefresh()
	if err != nil {
		return nil, err
	}
	_ = u.repo.SetHash(guid, hashedRefresh)
	/* if err != nil {
		return nil, err
	} */
	return tokenAll, nil
}
