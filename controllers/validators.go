package controllers

import (
	"BE/models"
	"errors"
	"strings"
)

func ValidateUserRegistration(user models.Users) error {
	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.Email) == "" {
		return errors.New("all fields (username, password, email) are required")
	}

	if len(user.Username) <= 5 {
		return errors.New("username must be at least 5 characters long")
	}

	if len(user.Password) <= 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func ValidateLoginCredentials(username, password string) error {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return errors.New("all fields (username & password) are required")
	}
	return nil
}

func ValidateArticle(post models.Posts) error {
	if strings.TrimSpace(post.Title) == "" || strings.TrimSpace(post.Content) == "" ||
		strings.TrimSpace(post.Category) == "" || strings.TrimSpace(post.Status) == "" {
		return errors.New("all fields (title, content, category, status) are required")
	}

	if len(post.Title) <= 10 {
		return errors.New("title must be at least 10 characters long")
	}

	if len(post.Content) <= 50 {
		return errors.New("content must be at least 50 characters long")
	}

	return nil
}
