package user

import (
	"context"

	"github.com/bsladewski/mojito/data"
)

// GetUserByID retrieves a user record by id.
func GetUserByID(ctx context.Context, id uint) (*User, error) {

	var item User

	if err := data.DB().Model(&User{}).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// GetUserByEmail retrieves a user record by email address.
func GetUserByEmail(ctx context.Context, email string) (*User, error) {

	var item User

	if err := data.DB().Model(&User{}).
		Where("LOWER(email) = LOWER(?)", email).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// SaveUser inserts or updates the supplied user record.
func SaveUser(ctx context.Context, item *User) error {
	return data.DB().Save(item).Error
}

// DeleteUser deletes the supplied user record.
func DeleteUser(ctx context.Context, item *User) error {
	return data.DB().Delete(item).Error
}
