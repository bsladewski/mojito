package auth

import (
	"time"

	"github.com/bsladewski/mojito/data"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// init migrates the database model.
func init() {
	data.DB().AutoMigrate(
		User{},
		UserAuth{},
	)

	// check if we should use mock data
	if !data.UseMockData() {
		return
	}

	// load mock data
	for _, user := range mockUsers {
		if err := data.DB().Create(&user).Error; err != nil {
			logrus.Fatal(err)
		}
	}
}

/* Data Types */

// User provides access to the mojito application.
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Email    string `gorm:"index,unique" json:"email"`
	Password string `json:"password"`
	Verified bool   `json:"verified"`

	LoggedOutAt time.Time `json:"logged_out_at"`
}

// UserAuth stores identifiers for validating user auth tokens.
type UserAuth struct {
	ID uint `gorm:"primarykey" json:"id"`

	UserID uint   `gorm:"index" json:"user_id"`
	UUID   string `gorm:"index" json:"uuid"`

	ExpiresAt time.Time `json:"expires_at"`
}

/* Mock Data */

// mockUsers defines mock data for the user type.
var mockUsers = []User{
	{
		ID:       1,
		Email:    "test@mojitobot.com",
		Password: "$2a$10$38cznnVvOXAd4fFZH/M89efgJP3LB0p2NnyXystHkRlxrSeL2tkvS", // mojito
		Verified: true,
	},
}
