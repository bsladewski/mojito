package auth

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

// GetUserAuthByID retrieves a user auth record by id.
func GetUserAuthByID(ctx context.Context,
	id uint) (*UserAuth, error) {

	var item UserAuth

	if err := data.DB().Model(&UserAuth{}).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// GetUserAuthByUUID retrieves a user auth record by UUID.
func GetUserAuthByUUID(ctx context.Context,
	uuid string) (*UserAuth, error) {

	var item UserAuth

	if err := data.DB().Model(&UserAuth{}).
		Where("uuid = ?", uuid).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// ListUserAuthByUserID retrieves all user auth records associated with the
// supplied user id.
func ListUserAuthByUserID(ctx context.Context,
	userID uint) ([]*UserAuth, error) {

	var items []*UserAuth

	if err := data.DB().Model(&UserAuth{}).
		Where("user_id = ?", userID).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}

// SaveUserAuth inserts or updates the supplied user auth record.
func SaveUserAuth(ctx context.Context, item *UserAuth) error {
	return data.DB().Save(item).Error
}

// DeleteUserAuth deletes the supplied user auth record.
func DeleteUserAuth(ctx context.Context, item *UserAuth) error {
	return data.DB().Delete(item).Error
}

// DeleteExpiredUserAuth deletes all expires user auth records associated with
// the specified user id.
func DeleteExpiredUserAuth(ctx context.Context, userID uint) error {
	return data.DB().
		Where("user_id = ?", userID).
		Where("expires_at < ?", time.Now()).
		Delete(&UserAuth{}).Error
}
