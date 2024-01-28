package main

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string `json:"password"`
}

func createUser(db *gorm.DB, user *User) error {
	hashPass, err :=
		bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPass)
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
