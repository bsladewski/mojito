package user

import (
	"context"
	"time"

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

// GetLoginByID retrieves a user login record by id.
func GetLoginByID(ctx context.Context,
	id uint) (*Login, error) {

	var item Login

	if err := data.DB().Model(&Login{}).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// GetLoginByUUID retrieves a user login record by UUID.
func GetLoginByUUID(ctx context.Context,
	uuid string) (*Login, error) {

	var item Login

	if err := data.DB().Model(&Login{}).
		Where("uuid = ?", uuid).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// ListLoginByUserID retrieves all user login records associated with the
// supplied user id.
func ListLoginByUserID(ctx context.Context,
	userID uint) ([]*Login, error) {

	var items []*Login

	if err := data.DB().Model(&Login{}).
		Where("user_id = ?", userID).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}

// SaveLogin inserts or updates the supplied user login record.
func SaveLogin(ctx context.Context, item *Login) error {
	return data.DB().Save(item).Error
}

// DeleteLogin deletes the supplied user login record.
func DeleteLogin(ctx context.Context, item *Login) error {
	return data.DB().Delete(item).Error
}

// DeleteExpiredLogin deletes all expires user login records associated with
// the specified user id.
func DeleteExpiredLogin(ctx context.Context, userID uint) error {
	return data.DB().
		Where("user_id = ?", userID).
		Where("expires_at < ?", time.Now()).
		Delete(&Login{}).Error
}
