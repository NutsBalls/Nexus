package services

import (
	"github.com/NutsBalls/Nexus/models"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// UserService handles business logic for users
type UserService struct {
	DB *gorm.DB
}

// CreateUser creates a new user and stores it in the database
func (us *UserService) CreateUser(username, email, password string) (models.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Username: username, Email: email, PasswordHash: string(hashedPassword), CreatedAt: time.Now()}
	if err := us.DB.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (us *UserService) Login(username, password string) (string, error) {
	var user models.User
	if err := us.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	// Here you can generate JWT token
	return "some-jwt-token", nil
}
