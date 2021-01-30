package feed

import (
	"time"

	"gorm.io/gorm"
)

/* Data Types */

// feedPlatform stores configuration needed to connect to a feed of market data.
type feedPlatform struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name     string        `gorm:"index" json:"name"`
	Enabled  bool          `json:"enabled"`
	BaseURL  string        `json:"base_url"`
	Interval time.Duration `json:"interval"`

	Securities []feedPlatformSecurity `json:"securities"`
}

// feedPlatformSecurity stores information about the price data we want to
// retrieve in a platform feed.
type feedPlatformSecurity struct {
	ID             uint `gorm:"primarykey" json:"id"`
	FeedPlatformID uint `json:"feed_platform_id"`

	Exchange          string `json:"exchange"`
	Ticker            string `json:"ticker"`
	ReferenceCurrency string `json:"reference_currency"`
}

/* Mock Data */

var mockFeedPlatforms = []feedPlatform{
	{
		ID:       1,
		Name:     PlatformCoinbase,
		Enabled:  true,
		BaseURL:  "wss://ws-feed.pro.coinbase.com",
		Interval: 60 * time.Second,
	},
}

var mockFeedPlatformSecurities = []feedPlatformSecurity{
	{
		ID:                1,
		FeedPlatformID:    1,
		Exchange:          exchangeCoinbase,
		Ticker:            "BTC",
		ReferenceCurrency: "USD",
	},
}
