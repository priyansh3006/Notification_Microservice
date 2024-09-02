package notification

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSms(phoneNumber string, message string) error {
	//func main() {
	// err := godotenv.Load("twilio.env")
	// if err != nil {
	// 	log.Fatalf("Can not load environment file")
	// }
	accountSid := "ACdb96e2cbacebea6e63327d4afe1950f5"
	authToken := "40359b6c9ef7468ca52369c670f9ee65"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo("+917023045653")  // Replace with phone number
	params.SetFrom("+17816946328") // Replace with Twilio phone number
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS:", err)
		return err
	}

	if resp.Sid != nil {
		fmt.Println("Message SID:", *resp.Sid)
	} else {
		fmt.Println("Message SID: unknown")
	}
	fmt.Println("Sent the message from the server using sms")
	return nil
}
