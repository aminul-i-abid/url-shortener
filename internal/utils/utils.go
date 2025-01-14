package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/aminul-i-abid/url-shortener/internal/models"
	"github.com/google/uuid"
)

func WriteJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := models.Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func GenerateShortCode() string {
	return uuid.New().String()[:6]
}

func ValidateURL(url string) error {
	if url == "" {
		return errors.New("URL cannot be empty")
	}

	urlRegex := regexp.MustCompile(`^(http|https)://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(url) {
		return errors.New("URL is invalid")
	}

	return nil
}
