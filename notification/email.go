package notification

import (
	"fmt"
	"log"

	gomail "gopkg.in/mail.v2"
)

func EmailSender(myMail string, userSubject string, userBody string) error {
	abc := gomail.NewMessage()

	abc.SetHeader("From", "priyansh3006@gmail.com")
	abc.SetHeader("To", "b22230@students.iitmandi.ac.in")
	abc.SetHeader("Subject", "Email")
	abc.SetBody("text/plain", "from my microservice")

	a := gomail.NewDialer("smtp.gmail.com", 587, "priyansh3006@gmail.com", "qfwq uyoq vszy psvu")

	if err := a.DialAndSend(abc); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sent the message from the server in email")
	return nil

}
