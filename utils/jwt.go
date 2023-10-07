package utils

import (
	"chicchat/models"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type userProfile struct {
	ID    uint
	Email string
}

func CreateAccessToken(email string, db *gorm.DB) (string, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("email %s not found", email)
		}
		return "", fmt.Errorf("error query data")
	}

	expTime := time.Now().Add(time.Minute * 60 * 3).Unix() // Default exp time should be 15 minute
	// Create the Claims
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": email,
		"type":  "access_token",
		"exp":   expTime,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	return t, err
}

func CreateRefreshToken(email string, db *gorm.DB) (string, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("email %s not found", email)
		}
		return "", fmt.Errorf("error query data")
	}

	expTime := time.Now().Add(time.Hour * 24 * 7).Unix()
	// Create the Claims
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": email,
		"type":  "refresh_token",
		"exp":   expTime,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET_KEY")))

	return t, err
}

func GetUserProfileFromResetPasswordToken(tokenString string) (userProfile, error) {
	var user userProfile
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_MAILER_SECRET_KEY")), nil
	})
	if err != nil {
		return user, fmt.Errorf("error verifying token or expired token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, fmt.Errorf("invalid claims format")
	}

	token_type, ok := claims["type"].(string)
	if !ok {
		return user, fmt.Errorf("type not found in token")
	}
	if token_type != "reset_password" {
		return user, fmt.Errorf("wrong token type")
	}

	userIdFloat, ok := claims["id"].(float64)
	if !ok {
		return user, fmt.Errorf("user ID not found in token")
	}
	user.ID = uint(userIdFloat)

	userEmail, ok := claims["email"].(string)
	if !ok {
		return user, fmt.Errorf("user email not found in token")
	}
	user.Email = userEmail

	return user, nil
}

func FormatBearerToken(token string) (string, error) {
	if token == "" || len(token) <= 7 || token[:7] != "Bearer " {
		return "", fmt.Errorf("invalid or expired access token, and no refresh token found")
	}
	return token[7:], nil
}

// You can return error without writing c.SendStatus with this function
func GetBearerToken(c *fiber.Ctx) (string, error) {
	token := c.Get("Authorization")

	if !strings.HasPrefix(token, "Bearer ") {
		return "", c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid token type"})
	}

	if len(token) <= 7 {
		return "", c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid or empty token"})
	}
	return token[7:], nil
}

func GetUserProfileFromToken(tokenString string) (userProfile, error) {
	var user userProfile
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return user, fmt.Errorf("error verifying token or expired token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, fmt.Errorf("")
	}

	userIdFloat, ok := claims["id"].(float64)
	if !ok {
		return user, fmt.Errorf("user ID not found in token")
	}
	user.ID = uint(userIdFloat)

	userEmail, ok := claims["email"].(string)
	if !ok {
		return user, fmt.Errorf("user email not found in token")
	}
	user.Email = userEmail

	return user, nil
}
