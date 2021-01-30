package email

import (
	"mojito/data"
	"mojito/env"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// init migrates the package model and loads the SMTP configuration.
func init() {

	data.DB().AutoMigrate(
		emailTemplate{},
		emailLog{},
	)

	// check if we should use mock data
	if !data.UseMockData() {
		return
	}

	// load mock data
	for _, t := range mockEmailTemplates {
		if err := data.DB().Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&t).Error; err != nil {
			logrus.Fatal(err)
		}
	}

	// retrieve SMTP settings from the environment
	smtpUsername = env.GetStringSafe(smtpUsernameVariable, "")
	smtpPassword = env.GetStringSafe(smtpPasswordVariable, "")
	smtpHost = env.GetStringSafe(smtpHostVariable, "")
	smtpPort = env.GetIntSafe(smtpPortVariable, 25)

	// retrieve SES settings from the environment.
	sesAccessKeyID = env.GetStringSafe(sesAccessKeyIDVariable, "")
	sesAccessKeySecret = env.GetStringSafe(sesAccessKeySecretVariable, "")
	sesRegion = env.GetStringSafe(sesRegionVariable, "")

	// determine email sending method based on environment configuration
	if smtpUsername != "" && smtpPassword != "" && smtpHost != "" {
		sendingMethod = sendingMethodSMTP
	} else if sesRegion != "" && sesAccessKeyID != "" && sesAccessKeySecret != "" {
		sendingMethod = sendingMehtodSES
	}

	// if no email sending method was configured log a fatal error
	if sendingMethod == "" {
		logrus.Fatal("no email sending method was specified")
	}

	logEmails = env.GetBoolSafe(logEmailsVariable, false)

	// retrieve default from and reply-to addresses
	defaultFromAddress = env.MustGetString(defaultFromAddressVariable)
	defaultReplyToAddress = env.MustGetString(defaultReplyToAddressVariable)

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
	// sesRegionVariable defines an environment variable for the AWS region to
	// use when sending emails.
	sesRegionVariable = "MOJITO_SES_REGION"
	// sesAccessKeyIDVariable defines an environment variable for the AWS access
	// key id to use when sending emails.
	sesAccessKeyIDVariable = "MOJITO_SES_ACCESS_KEY_ID"
	// sesAccessKeySecretVariable defines an environment variable for the AWS
	// access key secret to use when sending emails.
	sesAccessKeySecretVariable = "MOJITO_SES_ACCESS_KEY_SECRET"
	// defaultFromAddressVariable defines an environement variable for the
	// default email address used when sending emails.
	defaultFromAddressVariable = "MOJITO_DEFAULT_FROM_ADDRESS"
	// defaultReplyToAddressVariable defines an environment variable for the
	// default reply-to email address used when sending emails.
	defaultReplyToAddressVariable = "MOJITO_DEFAULT_REPLY_TO_ADDRESS"
	// logEmailsVariable defines an evironment variable that determines whether
	// we should log the results of sending emails.
	logEmailsVariable = "MOJITO_LOG_EMAILS"
)
