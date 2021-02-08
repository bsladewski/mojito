package feed

import (
	"mojito/market"
	"time"

	"gorm.io/gorm"
)

/* Data Types */

// platformFeed stores configuration needed to connect to a feed of market data.
type platformFeed struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	PlatformID uint            `json:"platform_id"`
	Platform   market.Platform `json:"platform"`

	Name     string        `gorm:"index" json:"name"`
	Enabled  bool          `json:"enabled"`
	BaseURL  string        `json:"base_url"`
	Interval time.Duration `json:"interval"`

	Securities []platformFeedSecurity `json:"securities"`
}

// platformFeedSecurity stores information about the price data we want to
// retrieve in a platform feed. Securities defined in this table will be
// subscribed to when the feed is initialized.
type platformFeedSecurity struct {
	ID             uint `gorm:"primarykey" json:"id"`
	PlatformFeedID uint `json:"platform_feed_id"`

	Exchange string `json:"exchange"`
	Ticker   string `json:"ticker"`
}

/* Mock Data */

var mockFeedPlatforms = []platformFeed{
	{
		ID:         1,
		PlatformID: 1,
		Enabled:    true,
		BaseURL:    "wss://ws-feed.pro.coinbase.com",
		Interval:   60 * time.Second,
	},
}

var mockFeedPlatformSecurities = []platformFeedSecurity{
	{
		ID:             1,
		PlatformFeedID: 1,
		Exchange:       exchangeCoinbase,
		Ticker:         "BTC",
	},
}
