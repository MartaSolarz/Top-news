package service

import (
	"errors"
	"regexp"
)

var ErrInvalidEmail = errors.New("invalid email")

type SaveEmailService struct {
	newsDB NewsDatabase
}

func NewSaveEmailService(newsDB NewsDatabase) *SaveEmailService {
	return &SaveEmailService{
		newsDB: newsDB,
	}
}

func (s *SaveEmailService) SaveEmail(email string) error {
	if !validateEmail(email) {
		return ErrInvalidEmail
	}

	err := s.newsDB.SaveEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
