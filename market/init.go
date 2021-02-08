package market

import (
	"mojito/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// init migrates the package model.
func init() {
	data.DB().AutoMigrate(
		Candlestick{},
		Platform{},
	)

	if !data.UseMockData() {
		return
	}

	// load mock data
	for _, p := range mockPlatforms {
		if err := data.DB().Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&p).Error; err != nil {
			logrus.Fatal(err)
		}
	}

	for _, e := range mockExchanges {
		if err := data.DB().Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&e).Error; err != nil {
			logrus.Fatal(err)
		}
	}

}
