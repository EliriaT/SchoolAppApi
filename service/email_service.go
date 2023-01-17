package service

import (
	"bytes"
	"fmt"
	"github.com/matcornic/hermes/v2"
	"log"
	"net/mail"
	"net/smtp"
)

const (
	SMTP_HOST = "smtp.gmail.com"
	SMTP_PORT = "587"
)

type EmailSender struct {
	emailServerLogin string

	auth   smtp.Auth
	hermes hermes.Hermes
}

// The object that sends email and authentificates to the email server
func NewEmailSender(login, pass string) *EmailSender {
	e := &EmailSender{}
	e.emailServerLogin = login
	e.auth = smtp.PlainAuth("", login, pass, SMTP_HOST)

	e.hermes = hermes.Hermes{
		Product: hermes.Product{
			Name:      "E-School",
			Link:      "https://github.com/EliriaT/SchoolAppApi",
			Logo:      "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
			Copyright: "Copyright © 2023 E-School. No rights reserved.",
		},
	}
	return e
}

// Creates an email for confirming account creation. A special link needs to be provided
func NewConfirmationEmail(finishAccountLink string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				"Your account has been created successfully",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To finish account creation, please click here:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Finish account creation",
						Link:  finishAccountLink,
					},
				},
			},
		},
	}
}

// Creates an email for resetting account password. A special link needs to be provided
func NewResetPasswordEmail(changePasswordLink string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				"A request for password reset has been sent",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To change your password, please click here:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Change Password",
						Link:  changePasswordLink,
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no action is required",
			},
		},
	}
}

// Sends an email. Needs a receiver which is an mail Address. Address contains the name of the receiver and their address. The title is the title of the email. The email itself is generated above.
func (e *EmailSender) SendEmail(receiver mail.Address, title string, email hermes.Email) error {
	email.Body.Name = receiver.Name

	htmlBody, _ := e.hermes.GenerateHTML(email)
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", title, mimeHeaders)))
	body.Write([]byte(htmlBody))

	receiverArray := []string{receiver.Address}
	err := smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, e.auth, e.emailServerLogin, receiverArray, body.Bytes())
	if err != nil {
		log.Println(err)
	}
	log.Println("Email Sent Successfully!")
	return err
}
