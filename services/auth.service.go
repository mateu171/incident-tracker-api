package services

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"incident-tracker-api/config"
	"incident-tracker-api/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailExists        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

func RegisterUser(email, password string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	var existingUser models.User

	err := config.DB.Where("email = ?", email).First(&existingUser).Error

	if err == nil {
		return nil, ErrEmailExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	var user models.User

	err := config.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.RegisteredClaims{
		Subject:   string(rune(user.ID)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	claims.Subject = formatUserID(user.ID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "development-secret-change-me"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func formatUserID(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
