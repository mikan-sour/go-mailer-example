package service

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/jedzeins/go-mailer/src/repository"
)

type ProfanityDetectionService interface {
	GetProfanities(context.Context, string) ([]*repository.BadWord, error)
}

type ProfanityDetectionServiceImpl struct {
	Repo            repository.ProfanityRepo
	ProfanitiesList []*repository.BadWord
}

func NewProfanityDetectionService(db *sql.DB) (*ProfanityDetectionServiceImpl, error) {
	profanityRepo, err := repository.NewRepo(context.Background(), db)
	if err != nil {
		return nil, err
	}
	return &ProfanityDetectionServiceImpl{Repo: profanityRepo}, nil
}

func (pds *ProfanityDetectionServiceImpl) GetProfanities(ctx context.Context, region string) error {
	words, err := pds.Repo.GetProfanities(ctx, region)
	if err != nil {
		return err
	}

	pds.ProfanitiesList = words

	return nil
}

func (pds *ProfanityDetectionServiceImpl) CheckAgainstProfanitiesList(text string) (string, error) {
	for _, word := range pds.ProfanitiesList {
		r, err := regexp.Compile(fmt.Sprintf("\\b%s(\\b|\\p{P})", word.Word))
		if err != nil {
			return "", err
		}
		if r.FindString(text) != "" {
			return word.Word, nil
		}
	}

	return "", nil
}
