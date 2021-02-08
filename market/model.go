package market

import (
	"time"

	"gorm.io/gorm"
)

/* Data Types */

// Candlestick stores price data for a specific ticker over an interval of time.
type Candlestick struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	OpensDay  bool `gorm:"index" json:"-"`
	OpensHour bool `gorm:"index" json:"-"`

	Exchange string  `gorm:"index" json:"exchange"`
	Ticker   string  `gorm:"index" json:"ticker"`
	Open     float64 `json:"open"`
	Close    float64 `json:"close"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Volume   int     `json:"volume"`
}

// Platform stores descriptive information about a platform that can be used for
// retrieving market data or as a brokerage.
type Platform struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Key         PlatformKey `gorm:"index" json:"key"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	SignupLink  string      `json:"signup_link"`

	HasPriceAPI bool `json:"has_price_api"`
}

// Exchange stores descriptive information about an exchange.
type Exchange struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	PlatformID uint     `json:"platform_id"`
	Platform   Platform `json:"platform"`

	Key  ExchangeKey `gorm:"index" json:"key"`
	Name string      `json:"name"`
}

/* Mock Data */

var mockPlatforms = []Platform{
	{
		ID:          1,
		Key:         PlatformCoinbase,
		Name:        "Coinbase",
		Description: "An exchange for trading a variety of cryptocurrencies.",
		SignupLink:  "https://www.coinbase.com/signup",
	},
	{
		ID:          2,
		Key:         PlatformAlpaca,
		Name:        "Alpaca",
		Description: "A commission-free stock trading API.",
		SignupLink:  "https://app.alpaca.markets/signup",
		HasPriceAPI: true,
	},
}

var mockExchanges = []Exchange{
	{
		ID:         1,
		PlatformID: 1,
		Key:        ExchangeCoinbase,
		Name:       "Coinbase",
	},
	{
		ID:         2,
		PlatformID: 2,
		Key:        ExchangeIEX,
		Name:       "IEX (Investors Exchange LLC)",
	},
	{
		ID:         3,
		PlatformID: 2,
		Key:        ExchangeNASDAQBX,
		Name:       "Nasdaq BX, Inc.",
	},
	{
		ID:         4,
		PlatformID: 2,
		Key:        ExchangeNASDAQPSX,
		Name:       "Nasdaq PSX",
	},
	{
		ID:         5,
		PlatformID: 2,
		Key:        ExchangeNYSENational,
		Name:       "NYSE National, Inc.",
	},
	{
		ID:         6,
		PlatformID: 2,
		Key:        ExchangeNYSEChicago,
		Name:       "NYSE Chicago, Inc.",
	},
}
