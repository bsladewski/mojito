package user

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
	)

	// check if we should use mock data
	if !data.UseMockData() {
		return
	}

	// load mock data
	for _, u := range mockUsers {
		if err := data.DB().Create(&u).Error; err != nil {
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
