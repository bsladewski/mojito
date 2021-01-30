package user

import (
	"mojito/data"
	"mojito/env"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// init migrates the package model and configures the user package. This
// function reads an access and refresh key from the environment for JWT
// signing, if these keys are not found the application will log a fatal error.
func init() {

	data.DB().AutoMigrate(
		User{},
		Login{},
	)

	// check if we should use mock data
	if !data.UseMockData() {
		return
	}

	// load mock data
	for _, u := range mockUsers {
		if err := data.DB().Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&u).Error; err != nil {
			logrus.Fatal(err)
		}
	}

	// get access key for signing access tokens
	accessKey = env.MustGetString(accessKeyVariable)

	// get refresh key for signing refresh tokens
	refreshKey = env.MustGetString(refreshKeyVariable)

	// configure access token expiration time
	accessExpirationHours = time.Duration(
		env.GetIntSafe(accessExpirationHoursVariable, 8)) * time.Hour

	// configure refresh token expiration time
	refreshExpirationHours = time.Duration(
		env.GetIntSafe(refreshExpirationHoursVariable, 168)) * time.Hour

}

const (
	// accessKeyVariable defines an environment variable for the key used to
	// sign JWT access tokens.
	accessKeyVariable = "MOJITO_ACCESS_KEY"
	// refreshKeyVariables defines an environment variable for the key used to
	// sign JWT refresh tokens.
	refreshKeyVariable = "MOJITO_REFRESH_KEY"
	// accessExpirationHoursVariable defines an environment variable for the
	// number of hours before we should consider an access token expired.
	accessExpirationHoursVariable = "MOJITO_ACCESS_EXPIRATION_HOURS"
	// refreshExpirationHoursVariable defines an environment variable for the
	// number of hours before we should consider a refresh token expired.
	refreshExpirationHoursVariable = "MOJITO_REFRESH_EXPIRATION_HOURS"
)
