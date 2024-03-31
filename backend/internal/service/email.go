package service

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"regexp"

	"gopkg.in/mail.v2"
)

var ErrInvalidEmail = errors.New("invalid email")

type EmailService struct {
	newsDB            NewsDatabase
	hostEmail         string
	hostEmailPassword string
	emailTopic        string
}

func NewSaveEmailService(newsDB NewsDatabase, hostEmail, hostEmailPass, emailTopic string) *EmailService {
	return &EmailService{
		newsDB:            newsDB,
		hostEmail:         hostEmail,
		hostEmailPassword: hostEmailPass,
		emailTopic:        emailTopic,
	}
}

func (s *EmailService) SaveEmail(email string) error {
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

func (s *EmailService) SendEmails() error {
	emails, err := s.newsDB.GetEmails()
	if err != nil {
		return err
	}

	log.Printf("Sending emails to %d subscribers\n", len(emails))

	message, err := s.createMessage()
	if err != nil {
		log.Printf("could not create message: %s", err.Error())
		return err
	}

	for _, email := range emails {
		err = s.sendEmail(email, message)
		if err != nil {
			log.Printf("could not send email to %s: %s\n", email, err.Error())
			return err
		}
	}

	return nil
}

func addOne(i int) int {
	return i + 1
}

func (s *EmailService) createMessage() (string, error) {
	articles, err := s.newsDB.GetNewsFromToday()
	if err != nil {
		log.Printf("could not get news: %s", err.Error())
		return "", err
	}

	tmpl, err := template.New("email.html").Funcs(template.FuncMap{"addOne": addOne}).ParseFiles("frontend/html/email.html")
	if err != nil {
		log.Printf("error parsing template: %s", err)
		return "", err
	}

	var message bytes.Buffer
	if err = tmpl.Execute(&message, articles); err != nil {
		log.Printf("error executing template: %s", err)
		return "", err
	}

	return message.String(), nil
}

func (s *EmailService) sendEmail(emailAddress, message string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.hostEmail)
	m.SetHeader("To", emailAddress)
	m.SetHeader("Subject", s.emailTopic)
	m.SetBody("text/html", message)

	d := mail.NewDialer("smtp.poczta.onet.pl", 465, s.hostEmail, s.hostEmailPassword)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
