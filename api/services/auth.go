package services

import (
	"backend-admin/api/models"
	"errors"
	"log"

	"firebase.google.com/go/auth"
	"gorm.io/gorm"
)

// AuthService provides authentication and authorization-related services.
type AuthService struct {
	FirebaseAuthClient *auth.Client
	DB                 *gorm.DB
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(firebaseAuthClient *auth.Client, db *gorm.DB) *AuthService {
	return &AuthService{
		FirebaseAuthClient: firebaseAuthClient,
		DB:                 db,
	}
}

func (s *AuthService) RegisterUser(uid, email string) error {
	log.Println("Start RegisterUser")

	// Check if the user is already registered in the database
	var existingUser models.Usuario
	if err := s.DB.First(&existingUser, "uid = ?", uid).Error; err == nil {
		log.Printf("User with UID %s is already registered\n", uid)
		return errors.New("User is already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error checking for existing user: %v\n", err)
		return err
	}
	log.Println("User is not registered, proceeding with insertion")

	// Insert the user into the database
	newUser := models.Usuario{
		ID:    uid,
		Email: email,
	}
	if err := s.DB.Create(&newUser).Error; err != nil {
		log.Printf("Error inserting user: %v\n", err)
		return err
	}

	log.Println("User registered successfully")
	return nil
}
