package market

import "mojito/data"

// init migrates the package model.
func init() {
	data.DB().AutoMigrate(
		Candlestick{},
	)
}
