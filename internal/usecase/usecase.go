package usecase

import (
	"log"

	"github.com/google/uuid"
	"github.com/paxaf/medodsTestEx/internal/repository"
	"github.com/paxaf/medodsTestEx/internal/tokens"
)

type UseCase interface {
	GetTokens(guid string, hashAgent string) (*tokens.Tokens, error)
	ValidateTokens(tokens tokens.Tokens) (string, string, bool)
	UpdateTokens(guid string) (*tokens.Tokens, error)
	ValidateJWT(tokens tokens.Tokens) (string, bool)
	UnathorizeUser(guid string) error
}

type usecase struct {
	repo repository.PgRepository
}

func NewUseCase(repo repository.PgRepository) UseCase {
	return &usecase{
		repo: repo,
	}
}
func (u *usecase) GetTokens(guid string, hashAgent string) (*tokens.Tokens, error) {
	err := uuid.Validate(guid)
	if err != nil {
		return nil, err
	}
	tokenAll := tokens.NewTokens(guid)
	hashedRefresh, err := tokenAll.GetHashedRefresh()
	if err != nil {
		return nil, err
	}
	err = u.repo.SetHash(guid, hashedRefresh, hashAgent)
	if err != nil {
		log.Println("ошибка sethash:", err)
		return nil, err
	}
	return tokenAll, nil
}

func (u *usecase) UpdateTokens(guid string) (*tokens.Tokens, error) {
	err := uuid.Validate(guid)
	if err != nil {
		return nil, err
	}
	tokenAll := tokens.NewTokens(guid)
	hashedRefresh, err := tokenAll.GetHashedRefresh()
	if err != nil {
		return nil, err
	}
	err = u.repo.UpdHash(guid, hashedRefresh)
	if err != nil {
		log.Println("ошибка sethash:", err)
		return nil, err
	}
	return tokenAll, nil
}

func (u *usecase) ValidateTokens(tokens tokens.Tokens) (string, string, bool) {
	guid, err := tokens.ValidateJWT()
	if err != nil {
		log.Println("JWT validate err", err)
		return "", "", false
	}
	hashInPg, hashedAgent, err := u.repo.GetHash(guid)
	if err != nil {
		log.Println("Get hash err", err)
		return "", "", false
	}
	return guid, hashedAgent, tokens.ValidateRefresh(hashInPg)
}

func (u *usecase) ValidateJWT(tokens tokens.Tokens) (string, bool) {
	guid, err := tokens.ValidateJWT()
	if err != nil {
		return "", false
	}
	return guid, true
}

func (u *usecase) UnathorizeUser(guid string) error {
	return u.repo.DeleteHash(guid)
}
