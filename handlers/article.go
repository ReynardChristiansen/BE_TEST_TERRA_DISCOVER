package handlers

import (
	"BE/config"
	"BE/models"
	"BE/controllers"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

var jwtKey = []byte("v5H8FhM2pLtJqTzXNVqz9fP5X4gq9T0zY+Q6Y8wHhLk=")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}


type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  "ERROR",
		Message: message,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := controllers.ValidateLoginCredentials(creds.Username, creds.Password); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.Users
	if err := config.DB.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			sendErrorResponse(w, "User not found", http.StatusUnauthorized)
		} else {
			sendErrorResponse(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		sendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		sendErrorResponse(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "SUCCESS", "token": tokenString})
}



func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user := models.Users{
		Username: creds.Username,
		Password: creds.Password,
		Email:    creds.Email,
	}

	if err := controllers.ValidateUserRegistration(user); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var existingUser models.Users
	if err := config.DB.Where("username = ? OR email = ?", creds.Username, creds.Email).First(&existingUser).Error; err == nil {
		sendErrorResponse(w, "Username or email already exists", http.StatusConflict)
		return
	} else if err != gorm.ErrRecordNotFound {
		sendErrorResponse(w, "Database error", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	if result := config.DB.Create(&user); result.Error != nil {
		sendErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "SUCCESS",
		Message: "User registered successfully",
	})
}



func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			sendErrorResponse(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(tokenStr, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			sendErrorResponse(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenParts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			sendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	var post models.Posts
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := controllers.ValidateArticle(post); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := config.DB.Create(&post); result.Error != nil {
		sendErrorResponse(w, "Failed to create article", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "SUCCESS",
		Message: "Article created successfully",
	})
}



func GetArticles(w http.ResponseWriter, r *http.Request) {
	var posts []models.Posts
	result := config.DB.Find(&posts)

	if result.Error != nil {
		sendErrorResponse(w, "Failed to fetch articles", http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		sendErrorResponse(w, "No articles found", http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
