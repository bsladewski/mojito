package feed

import (
	"context"
	"mojito/data"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// init migrates the package model and connects to any configured platform
// feeds.
func init() {

	// migrate the package model
	data.DB().AutoMigrate(
		feedPlatform{},
		feedPlatformSecurity{},
	)

	// load mock data if the server is configured to use it
	if data.UseMockData() {

		for _, p := range mockFeedPlatforms {
			if err := data.DB().Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&p).Error; err != nil {
				logrus.Fatal(err)
			}
		}

		for _, s := range mockFeedPlatformSecurities {
			if err := data.DB().Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&s).Error; err != nil {
				logrus.Fatal(err)
			}
		}

	}

	// retrieve all enabled platforms
	platforms, err := ListPlatform(context.Background(), data.DB(),
		ptrToBool(true))
	if err != nil {
		logrus.Fatal(err)
	}

	// initialize the feed for each platform
	for _, platform := range platforms {
		if _, err := Connect(platform); err != nil {
			logrus.Fatal(err)
		}
	}

}
