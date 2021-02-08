package feed

import (
	"context"

	"gorm.io/gorm"
)

// ListPlatform retrieves all feed platforms, takes an optional flag that can be
// used to filter by enabled platforms.
func ListPlatform(ctx context.Context, db *gorm.DB,
	enabled *bool) ([]*platformFeed, error) {

	var items []*platformFeed

	res := db.Preload("Platform").Preload("Securities").Model(&platformFeed{})

	if enabled != nil {
		res = res.Where("enabled = ?", *enabled)
	}

	if err := res.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}
