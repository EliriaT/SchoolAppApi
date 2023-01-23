package service

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"text/template"
)

const (
	SMTP_HOST = "smtp.gmail.com"
	SMTP_PORT = "587"
)

var ErrSendingEmail = errors.New("Error Sending Email")

type EmailService interface {
	SendRegisterEmail(name, school, role, userPassword, QRcode, receiver, link string) error
	SendChangePasswordEmail(link string, receiver string) error
}

type emailService struct {
	from string
	auth smtp.Auth
}

// The object that sends email and authentificates to the email server
func NewEmailService(login, pass string) EmailService {
	e := emailService{}
	e.from = login
	e.auth = smtp.PlainAuth("", login, pass, SMTP_HOST)

	return e
}

// subject string, templatePath string, to []string
func (e emailService) SendRegisterEmail(name, school, role, userPassword, QRcode, receiver, link string) error {

	var body bytes.Buffer
	t, err := template.ParseFiles("service/email-template/email-template.html")
	if err != nil {
		return err
	}

	err = t.Execute(&body, struct {
		Name     string
		School   string
		Role     string
		Email    string
		Password string
		QRCode   string
		Link     string
	}{
		Name:     name,
		School:   school,
		Role:     role,
		Email:    receiver,
		Password: userPassword,
		QRCode:   QRcode,
		Link:     link,
	})
	if err != nil {
		return err
	}

	err = smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, e.auth, e.from, []string{receiver}, body.Bytes())

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
func (e emailService) SendChangePasswordEmail(link string, receiver string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles("service/email-template/recover-password-template.html")
	if err != nil {
		return err
	}

	err = t.Execute(&body, struct {
		Link string
	}{
		Link: link,
	})
	if err != nil {
		return err
	}

	err = smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, e.auth, e.from, []string{receiver}, body.Bytes())

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
