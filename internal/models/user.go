package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Phone    string
	Orders   []Order `gorm:"foreignKey:UserID"`
}

func CreateUser(db *gorm.DB, user *User) error {
	// Mã hóa mật khẩu trước khi lưu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return db.Create(user).Error
}

func FindByUsername(db *gorm.DB, name string) (*User, error) {
	var user User
	err := db.Where("name = ?", name).First(&user).Error
	return &user, err
}
