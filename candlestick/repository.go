package candlestick

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// GetLastByTicker retrieves the most recently added candlestick for the
// specified exchange and ticker.
func GetLastByTicker(ctx context.Context, db *gorm.DB, exchange,
	ticker string) (Candlestick, error) {

	var item Candlestick

	if err := db.Model(&Candlestick{}).
		Where("exchange = ? AND ticker = ?", exchange, ticker).
		Last(&item).Error; err != nil {
		return Candlestick{}, err
	}

	return item, nil

}

// ListByTicker retrieves all candlestick records associated with the specified
// ticker and between the specified start and end date.
func ListByTicker(ctx context.Context, db *gorm.DB, exchange, ticker string,
	hourly, daily bool, startDate, endDate time.Time) ([]Candlestick, error) {

	var items []Candlestick

	res := db.Model(&Candlestick{}).
		Where("exchange = ? AND ticker = ?", exchange, ticker).
		Where("created_at > ? AND created_at < ?", startDate, endDate)

	if hourly {
		res = res.Where("opens_hour")
	}

	if daily {
		res = res.Where("opens_day")
	}

	if err := res.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}

// ListExchanges retrieves all exchanges for which candlestick data exists.
func ListExchanges(ctx context.Context, db *gorm.DB) ([]string, error) {

	var exchanges []string

	rows, err := db.Raw("SELECT DISTINCT exchange FROM candlesticks").Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var exchange sql.NullString

		if err := rows.Scan(&exchange); err != nil {
			return nil, err
		}

		exchanges = append(exchanges, exchange.String)
	}

	return exchanges, nil

}

// ListTickers retrieves all tickers for which candlestick data exists.
func ListTickers(ctx context.Context, db *gorm.DB,
	exchange string) ([]string, error) {

	var tickers []string

	rows, err := db.Raw(
		"SELECT DISTINCT ticker FROM candlesticks WHERE exchange = ?",
		exchange).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ticker sql.NullString

		if err := rows.Scan(&ticker); err != nil {
			return nil, err
		}

		tickers = append(tickers, ticker.String)
	}

	return tickers, nil

}

// Save inserts or updates the supplied candlestick record.
func Save(ctx context.Context, db *gorm.DB, item Candlestick) error {
	return db.Save(&item).Error
}
