package email

import (
	"stori-service/src/libs/dto"
	"stori-service/src/libs/env"

	"github.com/go-gomail/gomail"
)

var (
	emailServer   = env.EmailServer
	emailAcount   = env.EmailAccount
	emailPort     = env.EmailPort
	emailPassword = env.EmailPassword
)

func SendEmail(movementList *dto.MovementList) error {
	m := gomail.NewMessage()
	m.SetHeader("From", emailAcount)
	m.SetHeader("To", movementList.Customer.Email)
	m.SetHeader("Subject", "Balance")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	d := gomail.NewDialer(emailServer, emailPort, emailAcount, emailPassword)

	return d.DialAndSend(m)
}
