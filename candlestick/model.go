package candlestick

import (
	"time"

	"mojito/data"
)

// init migrates the database model.
func init() {
	data.DB().AutoMigrate(
		Candlestick{},
	)
}

/* Data Types */

// Candlestick stores price data for a specific ticker over an interval of time.
type Candlestick struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	Exchange string  `gorm:"index" json:"exchange"`
	Ticker   string  `gorm:"index" json:"ticker"`
	Open     float64 `json:"open"`
	Close    float64 `json:"close"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Volume   int     `json:"volume"`
}
