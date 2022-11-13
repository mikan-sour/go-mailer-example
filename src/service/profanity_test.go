package service

import (
	"testing"

	"github.com/jedzeins/go-mailer/src/repository"
)

func TestCheckAgainstProfanitiesList(t *testing.T) {
	service := &ProfanityDetectionServiceImpl{}

	service.ProfanitiesList = []*repository.BadWord{{Word: "assface"}}

	text := "he has an assface"

	word, err := service.CheckAgainstProfanitiesList(text)
	if err != nil {
		t.Fatalf(`got %s, expected nil `, err.Error())
	}

	if word == "" {
		t.Fatalf(`got nothing but expected "test" `)
	}

	word, err = service.CheckAgainstProfanitiesList("no profanity here...")
	if err != nil {
		t.Fatalf(`got %s, expected nil `, err.Error())
	}

	if word != "" {
		t.Fatalf(`got %s but expected empty string`, word)
	}

}
