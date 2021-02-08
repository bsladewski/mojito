package user

import (
	"time"

	"gorm.io/gorm"
)

/* Data Types */

// User provides access to the mojito application.
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Email    string `gorm:"index,unique" json:"email"`
	Password string `json:"password"`

	Admin     bool   `json:"admin"`      // admins have the broadest set of user permissions
	SecretKey string `json:"secret_key"` // used to sign tokens when generating links for this user
	Verified  bool   `json:"verified"`   // whether the user has completed email verification

	LoggedOutAt *time.Time `json:"logged_out_at"` // records the last time the user explicitly logged out

	Settings UserSettings
}

// UserSettings stores settings for a user account.
type UserSettings struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint `gorm:"index" json:"user_id"`

	// Coinbase credentials for executing trades via API
	CoinbaseAPIKey    string `json:"coinbase_api_key"`
	CoinbaseSignature string `json:"coinbase_signature"`

	// Alpaca credentials for executing trades via API
	AlpacaAPIKey    string `json:"alpaca_api_key"`
	AlpacaSecretKey string `json:"alpaca_secret_key"`
}

// Login stores identifiers for validating user auth tokens.
type Login struct {
	ID uint `gorm:"primarykey" json:"id"`

	UserID uint   `gorm:"index" json:"user_id"`
	UUID   string `gorm:"index" json:"uuid"` // uniquely identifies a refresh token

	ExpiresAt time.Time `json:"expires_at"` // records when a refresh token will expire
}

/* Mock Data */

var mockUsers = []User{
	{
		ID:        1,
		Email:     "admin@mojitobot.com",
		Password:  "$2a$10$38cznnVvOXAd4fFZH/M89efgJP3LB0p2NnyXystHkRlxrSeL2tkvS", // mojito
		Admin:     true,
		SecretKey: "8bf83c80-f235-461e-9bd7-00c83a5cfff8",
		Verified:  true,
	},
	{
		ID:        2,
		Email:     "test@mojitobot.com",
		Password:  "$2a$10$rX27aiSnPB1pSSez49kJDe2EOzih77M1nbGfL7cmd5Aw8FM2asY3m", // mojito
		SecretKey: "43ee0e83-dc81-4263-8bb0-6ccddff8586d",
		Verified:  true,
	},
}
