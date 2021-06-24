package twillo

import (
	"../../config"

	"github.com/sfreiberg/gotwilio"
)

// SendVerifySMS sends sms to phone number via twillo
func SendVerifySMS(phone, verifyCode string) error {
	accountSid := config.TwilloSid
	authToken := config.TwilloToken
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := config.TwilloFrom
	to := phone
	message := "Your GoGo verification code is " + verifyCode
	_, _, err := twilio.SendSMS(from, to, message, "", "")
	return err
}
