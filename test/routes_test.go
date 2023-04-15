package test

import (
	"encoding/json"
	"fmt"
	"go_jwt/controllers"
	"go_jwt/database"
	"go_jwt/models"
	"go_jwt/routes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func createTestUser(db *gorm.DB) *models.User {
	password, _ := bcrypt.GenerateFromPassword([]byte("test-password"), 14)
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: password,
	}
	db.Create(&user)
	return &user
}

func getTokenString(user *models.User) string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    fmt.Sprint(user.Id),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	token, _ := claims.SignedString([]byte(controllers.SecretKey))

	return token
}

func TestRegister(t *testing.T) {
	// Setup
	db := database.DB

	app := fiber.New()
	routes.Setup(app)

	reqBody := `{"name": "Test User", "email": "test@example.com", "password": "test-password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that the user was created in the DB
	var user models.User
	result := db.First(&user, "email = ?", "test@example.com")
	assert.Nil(t, result.Error)
	assert.Equal(t, "Test User", user.Name)
}

func TestLogin(t *testing.T) {
	// Setup
	db := database.DB

	app := fiber.New()
	routes.Setup(app)

	user := createTestUser(db)

	reqBody := `{"email": "test@example.com", "password": "test-password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that a JWT token was set in the response cookie
	cookie := resp.Cookies()[0]
	assert.Equal(t, "jwt", cookie.Name)
	assert.NotEmpty(t, cookie.Value)

	// Decode the token to check its content
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(controllers.SecretKey), nil
	})
	assert.Nil(t, err)
	assert.True(t, token.Valid)

	claims := token.Claims.(*jwt.StandardClaims)
	assert.Equal(t, fmt.Sprint(user.Id), claims.Issuer)
}

func TestUser(t *testing.T) {
	// create a new test fiber app
	app := fiber.New()

	// setup database and migrate User model
	database.Connect()
	database.DB.AutoMigrate(&models.User{})

	// create a test user
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	database.DB.Create(&user)

	// generate a JWT token for the user
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    fmt.Sprintf("%d", user.Id),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})
	token, _ := claims.SignedString([]byte(controllers.SecretKey))

	// create a new HTTP GET request with the JWT cookie
	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: token})

	// execute the request and get the response
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// decode the response body into a User struct
	var responseBody models.User
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.NotNil(t, responseBody)

	// check the user's name and email in the response body
	assert.Equal(t, user.Name, responseBody.Name)
	assert.Equal(t, user.Email, responseBody.Email)
}
