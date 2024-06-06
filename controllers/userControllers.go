package controllers

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dennisenwiya/Go-AUTH/initializers"
	"github.com/dennisenwiya/Go-AUTH/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Claims struct for JWT
var jwtKey = []byte(initializers.SecretKey)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// UserRegister handles user registration
func UserRegister(c *gin.Context) {
	var userInput struct {
		Firstname string `json:"firstname" binding:"required"`
		Lastname  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Privilege string `json:"privilege" binding:"required"`
	}

	// c.Bind(&userInput)
	// uData := models.User{Firstname: userInput.Firstname, Lastname: userInput.Lastname, Email: userInput.Email,
	// 	Username: userInput.Username, Password: userInput.Password, Privilege: userInput.Privilege,
	// }
	// fmt.Print(uData)
	// Bind and validate input
	if err := c.ShouldBindJSON(&userInput); err != nil {
		if err == io.EOF {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body cannot be empty"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		log.Printf("Failed to bind JSON: %v", err)
		return
	}

	// // Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		log.Printf("Failed to hash password: %v", err)
		return
	}

	// // Create a User model with hashed password
	user := models.User{
		Firstname: userInput.Firstname,
		Lastname:  userInput.Lastname,
		Email:     userInput.Email,
		Username:  userInput.Username,
		Password:  string(hashedPassword),
		Privilege: userInput.Privilege,
	}

	// // Save user to database
	if result := initializers.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		log.Printf("Failed to save user to database: %v", result.Error)
		return
	}

	// Exclude the password from the response
	userResponse := models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Username:  user.Username,
		Privilege: user.Privilege,
	}

	// Return the created user
	c.JSON(http.StatusOK, gin.H{"user": userResponse})
}

// Placeholder for user login function
// UserLogin handles user login
func UserLogin(c *gin.Context) {
	var loginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind and validate input
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by username
	var user models.User
	if err := initializers.DB.Where("username = ?", loginInput.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		}
		return
	}

	// Compare the hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: loginInput.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token and user details (excluding the password)
	userResponse := models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Username:  user.Username,
		Privilege: user.Privilege,
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  userResponse,
		"token": tokenString,
	})
}

// Placeholder for user update function
func UserUpdate(c *gin.Context) {
	// Implement user update logic here
}

// Placeholder for user delete function
func UserDelete(c *gin.Context) {
	// Implement user delete logic here
}
