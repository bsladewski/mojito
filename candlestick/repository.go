package candlestick

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// ListByTicker retrieves all candlestick records associated with the specified
// ticker and between the specified start and end date.
func ListByTicker(ctx context.Context, db *gorm.DB, exchange, ticker string,
	startDate, endDate time.Time) ([]Candlestick, error) {

	var items []Candlestick

	if err := db.Model(&Candlestick{}).
		Where("exchange = ? AND ticker = ?", exchange, ticker).
		Where("created_at > ? AND created_at < ?", startDate, endDate).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}

// Save inserts or updates the supplied candlestick record.
func Save(ctx context.Context, db *gorm.DB, item Candlestick) error {
	return db.Save(&item).Error
}
