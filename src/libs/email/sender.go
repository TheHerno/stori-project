package email

import (
	"fmt"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/env"
	"stori-service/src/utils/constant"

	"github.com/go-gomail/gomail"
)

var (
	emailServer   = env.EmailServer
	emailAcount   = env.EmailAccount
	emailPort     = env.EmailPort
	emailPassword = env.EmailPassword
)

// getTransactionByMonth returns a map of transactions by month
func getTransactionByMonth(movements []entity.Movement) string {
	transactionByMonth := make(map[string]int)
	for _, movement := range movements {
		transactionByMonth[movement.Date.Format("January")]++
	}
	list := ""
	for month, count := range transactionByMonth {
		list += fmt.Sprintf("Number of transactions in %s: %d<br>", month, count)
	}
	return list
}

// getAvgDebit is a partial application of getAvgDebitOrCredit with debit as the movement type
func getAvgDebit(movements []entity.Movement) float64 {
	return getAvgDebitOrCredit(movements, constant.OutcomeType)
}

// getAvgCredit is a partial application of getAvgDebitOrCredit with credit as the movement type
func getAvgCredit(movements []entity.Movement) float64 {
	return getAvgDebitOrCredit(movements, constant.IncomeType)
}

// getAvgDebitOrCredit returns the average debit or credit amount
func getAvgDebitOrCredit(movements []entity.Movement, movementType int) float64 {
	var sum float64
	var count int
	for _, movement := range movements {
		if movement.Type == movementType {
			sum += movement.Quantity
			count++
		}
	}
	return sum / float64(count)
}

func getHTML(movementList *dto.MovementList) string {
	listByMonth := getTransactionByMonth(movementList.Movements)
	avgDebit := getAvgDebit(movementList.Movements)
	avgCredit := getAvgCredit(movementList.Movements)
	currentBalance := movementList.Movements[len(movementList.Movements)-1].Available
	storiLogoURL := "https://dd7tel2830j4w.cloudfront.net/f1650918197627x637468688019988200/Stori%20splash.svg"
	return fmt.Sprintf(`
		<center>
			<img src="%s" alt="stori logo"></img>
		</center>
		<p>
		Hello, <strong>%s</strong>!<br>
		</p>
		<p>
		Your total balance is: <strong>%.2f</strong>
		</p>
		<p>
		%s
		</p>
		<p>
		Average debit amount: <strong>%.2f</strong>
		Average credit amount: <strong>%.2f</strong>
		</p>
	`, storiLogoURL, movementList.Customer.Name, currentBalance, listByMonth, avgDebit, avgCredit)
}

func SendEmail(movementList *dto.MovementList) error {
	m := gomail.NewMessage()
	m.SetHeader("From", emailAcount)
	m.SetHeader("To", movementList.Customer.Email)
	m.SetHeader("Subject", "Balance")
	m.SetBody("text/html", getHTML(movementList))

	d := gomail.NewDialer(emailServer, emailPort, emailAcount, emailPassword)

	return d.DialAndSend(m)
}
