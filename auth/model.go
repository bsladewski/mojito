package auth

import (
	"time"

	"github.com/bsladewski/mojito/data"
)

// init migrates the database model.
func init() {
	data.DB().AutoMigrate(
		Login{},
	)
}

/* Data Types */

// Login stores identifiers for validating user auth tokens.
type Login struct {
	ID uint `gorm:"primarykey" json:"id"`

	UserID uint   `gorm:"index" json:"user_id"`
	UUID   string `gorm:"index" json:"uuid"`

	ExpiresAt time.Time `json:"expires_at"`
}
