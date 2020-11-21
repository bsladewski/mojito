package email

import (
	"fmt"

	"github.com/bsladewski/mojito/env"
	"gorm.io/gorm"

	"gopkg.in/gomail.v2"
)

// init loads the SMTP configuration.
func init() {

	// retrieve SMTP settings from the environment
	smtpUsername = env.MustGetString(smtpUsernameVariable)
	smtpPassword = env.MustGetString(smtpPasswordVariable)
	smtpHost = env.MustGetString(smtpHostVariable)
	smtpPort = env.GetIntSafe(smtpPortVariable, 25)

	// retrieve default from address
	defaultFromAddress = env.MustGetString(defaultFromAddressVariable)

}

const (
	// smtpUsernameVariable defines an environment variable for the SMTP
	// username.
	smtpUsernameVariable = "MOJITO_SMTP_USERNAME"
	// smtpPasswordVariable defines an environment variable for the SMTP
	// password.
	smtpPasswordVariable = "MOJITO_SMTP_PASSWORD"
	// smtpHostVariable defines an evironment variable for the SMTP host.
	smtpHostVariable = "MOJITO_SMTP_HOST"
	// smtpPortVariable defines an environment variable for the SMTP port.
	smtpPortVariable = "MOJITO_SMTP_PORT"
	// defaultFromAddressVariable defines an environement variable for the
	// default email address used when sending emails.
	defaultFromAddressVariable = "MOJITO_DEFAULT_FROM_ADDRESS"
)

// smtpUsername is used to authenticate with an SMTP server to send emails.
var smtpUsername string

// smtpPassword is used to authenticate with an SMTP server to send emails.
var smtpPassword string

// smtpHost is the host of an SMTP server to use for sending emails.
var smtpHost string

// smtpPort is the port of an SMTP server to use for sending emails.
var smtpPort int

// smtpEnabled stores whether we able to get the SMTP configuration from the
// environment.
var smtpEnabled bool

// defaultFromAddress stores the default application from email address.
var defaultFromAddress string

// DefaultFromAddress is the application default from email address.
func DefaultFromAddress() string {
	return defaultFromAddress
}

// SendEmailTemplate formats the specified email template and sends the email
// through SMTP.
func SendEmailTemplate(
	from string,
	to, cc, bcc []string,
	templateTitle TemplateTitle,
	data interface{},
) error {

	// execute the email template
	subject, bodyText, bodyHTML, err := ExecuteTemplate(templateTitle, data)
	if err != nil {
		return err
	}

	// wrap HTML email body with header and footer
	_, _, newBodyHTML, err := ExecuteTemplate(templateTitleHeaderFooter,
		struct{ Body string }{bodyHTML})
	if err != nil && err == gorm.ErrRecordNotFound {
		return err
	} else if err == nil {
		bodyHTML = newBodyHTML
	}

	// send the email
	return SendEmail(from, to, cc, bcc, subject, bodyText, bodyHTML)

}

// SendEmail sends an email through SMTP.
func SendEmail(
	from string,
	to, cc, bcc []string,
	subject, bodyText, bodyHTML string,
) error {

	fmt.Println(from, to, cc, bcc, subject, bodyText, bodyHTML)

	// initialize SMTP client
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// build email message
	message := gomail.NewMessage()

	// set sender
	message.SetHeader("From", from)

	// set recipients
	message.SetHeader("To", to...)

	if len(cc) > 0 {
		message.SetHeader("Cc", cc...)
	}

	if len(bcc) > 0 {
		message.SetHeader("Bcc", bcc...)
	}

	// set subject
	if subject != "" {
		message.SetHeader("Subject", subject)
	}

	// set contents
	if bodyText != "" {
		message.SetBody("text/plain", bodyText)
	}

	if bodyHTML != "" {
		message.SetBody("text/html", bodyHTML)
	}

	// send email
	return dialer.DialAndSend(message)

}
