package user

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// GetUserByID retrieves a user record by id.
func GetUserByID(ctx context.Context, db *gorm.DB, id uint) (*User, error) {

	var item User

	if err := db.Model(&User{}).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// GetUserByEmail retrieves a user record by email address.
func GetUserByEmail(ctx context.Context, db *gorm.DB,
	email string) (*User, error) {

	var item User

	if err := db.Model(&User{}).
		Where("LOWER(email) = LOWER(?)", email).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// SaveUser inserts or updates the supplied user record.
func SaveUser(ctx context.Context, db *gorm.DB, item *User) error {
	return db.Save(item).Error
}

// DeleteUser deletes the supplied user record.
func DeleteUser(ctx context.Context, db *gorm.DB, item *User) error {
	return db.Delete(item).Error
}

// GetLoginByID retrieves a user login record by id.
func GetLoginByID(ctx context.Context, db *gorm.DB,
	id uint) (*Login, error) {

	var item Login

	if err := db.Model(&Login{}).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// GetLoginByUUID retrieves a user login record by UUID.
func GetLoginByUUID(ctx context.Context, db *gorm.DB,
	uuid string) (*Login, error) {

	var item Login

	if err := db.Model(&Login{}).
		Where("uuid = ?", uuid).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}

// ListLoginByUserID retrieves all user login records associated with the
// supplied user id.
func ListLoginByUserID(ctx context.Context, db *gorm.DB,
	userID uint) ([]*Login, error) {

	var items []*Login

	if err := db.Model(&Login{}).
		Where("user_id = ?", userID).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}

// SaveLogin inserts or updates the supplied user login record.
func SaveLogin(ctx context.Context, db *gorm.DB, item *Login) error {
	return db.Save(item).Error
}

// DeleteLogin deletes the supplied user login record.
func DeleteLogin(ctx context.Context, db *gorm.DB, item *Login) error {
	return db.Delete(item).Error
}

// DeleteExpiredLogin deletes all expires user login records associated with
// the specified user id.
func DeleteExpiredLogin(ctx context.Context, db *gorm.DB, userID uint) error {
	return db.
		Where("user_id = ?", userID).
		Where("expires_at < ?", time.Now()).
		Delete(&Login{}).Error
}
