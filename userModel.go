package main

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	// fmt.Println("Password: ", hashPass)
	user.Password = string(hashPass)
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func LoginUser(db *gorm.DB, user *User) (string, error) {
	selectedUser := new(User)
	result := db.Where("email =?", user.Email).First(selectedUser)
	if result.Error != nil {
		return "", result.Error
	}
	// compare hashed password
	err := bcrypt.CompareHashAndPassword([]byte(selectedUser.Password),
		[]byte(user.Password))
	if err != nil {
		return "", err
	}
	// generate token
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": selectedUser.Email,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		})
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
