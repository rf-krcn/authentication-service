package helpers

import (
	"errors"

	"gorm.io/gorm"
)

var db *gorm.DB

func New(dbPool *gorm.DB) {
	db = dbPool
	db.AutoMigrate(&User{})
}

// GetAll returns a slice of all users, sorted by last name
func GetAll() ([]*User, error) {

	var users []*User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// GetByEmail returns one user by email
func GetByUserName(userName string) (*User, error) {

	var user User
	result := db.Where("user_name = ?", userName).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *User) DeleteByID(id int) error {

	db := db.Delete(&User{}, id)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func Insert(user User) error {

	// Check if the username is already taken
	var existingUser User
	if err := db.Where("user_name = ?", user.UserName).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	result := db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ResetPasswordBD(password string, id uint) error {

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	result := db.Model(&User{}).Where("id = ?", id).Update("password", hashedPassword)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
