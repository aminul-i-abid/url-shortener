package services

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aminul-i-abid/url-shortener/internal/db"
	"github.com/aminul-i-abid/url-shortener/internal/models"
	"github.com/aminul-i-abid/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	_ "github.com/aminul-i-abid/url-shortener/docs"
)

// @title URL Shortener API
// @version 1.0
// @description API for managing short URLs.
// @contact.name Abid
// @contact.email aminul.i.abid@gmail.com
// @host localhost:8080
// @BasePath /api/v1

func Routes(r *gin.RouterGroup) {
	// Health check route
	r.GET("/health", handleHealthCheck)

	// Shorten URL routes
	r.POST("/shorten", handleCreateShortLink)
	r.GET("/shorten/:shortCode", handleRedirect)
	r.PUT("/shorten/:shortCode", handleUpdateShortLink)
	r.DELETE("/shorten/:shortCode", handleDeleteShortLink)
	r.GET("/shorten/:shortCode/stats", handleShortLinkStats)
}

// @Summary Check server health
// @Description Health check endpoint to verify if the server is running
// @Tags Health
// @Success 200 {string} string "Server is healthy"
// @Router /api/v1/health [get]
// Health check handler
func handleHealthCheck(c *gin.Context) {
	utils.WriteJSON(c.Writer, http.StatusOK, "OK", "Server is healthy")
}

// @Summary Create a short URL
// @Description Create a new short URL from a long URL
// @Tags Shorten
// @Accept json
// @Produce json
// @Param payload body models.ShortURL true "Request body"
// @Success 201 {object} models.ShortURL
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/shorten [post]
func handleCreateShortLink(c *gin.Context) {
	var req models.ShortURL
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}
	if err := utils.ValidateURL(req.URL); err != nil {
		utils.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid URL", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shortURL := &models.ShortURL{
		ID:          primitive.NewObjectID(),
		URL:         req.URL,
		ShortCode:   utils.GenerateShortCode(),
		CreatedAt:   time.Now().Format(time.RFC3339),
		AccessCount: 0,
	}

	_, err := db.Collection().InsertOne(ctx, shortURL)
	if err != nil {
		log.Printf("Database error: %v", err)
		utils.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to create short URL", nil)
		return
	}

	utils.WriteJSON(c.Writer, http.StatusCreated, "Short link created successfully", shortURL)
}

// @Summary Redirect to the original URL
// @Description Redirect a short code to the original URL
// @Tags Shorten
// @Param shortCode path string true "Short Code"
// @Success 302 {string} string "Redirect to original URL"
// @Failure 404 {object} models.Response
// @Router /api/v1/shorten/{shortCode} [get]
func handleRedirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		utils.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid short code", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shortURL models.ShortURL
	err := db.Collection().FindOne(ctx, bson.M{"shortCode": shortCode}).Decode(&shortURL)
	if err != nil {
		utils.WriteJSON(c.Writer, http.StatusNotFound, "Short URL not found", nil)
		return
	}

	shortURL.AccessCount++
	_, err = db.Collection().UpdateOne(ctx, bson.M{"shortCode": shortCode}, bson.M{"$set": bson.M{"accessCount": shortURL.AccessCount}})
	if err != nil {
		log.Printf("Failed to update access count: %v", err)
	}

	c.Redirect(http.StatusTemporaryRedirect, shortURL.URL)
}

// @Summary Update a short URL
// @Description Update the original URL for a given short code
// @Tags Shorten
// @Accept json
// @Produce json
// @Param shortCode path string true "Short Code"
// @Param payload body models.ShortURL true "Request body"
// @Success 200 {object} models.ShortURL
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/shorten/{shortCode} [put]
func handleUpdateShortLink(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		utils.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid short code", nil)
		return
	}

	var req models.ShortURL
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}
	if err := utils.ValidateURL(req.URL); err != nil {
		utils.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid URL", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shortURL models.ShortURL
	err := db.Collection().FindOne(ctx, bson.M{"shortCode": shortCode}).Decode(&shortURL)
	if err != nil {
		utils.WriteJSON(c.Writer, http.StatusNotFound, "Short URL not found", nil)
		return
	}

	shortURL.URL = req.URL
	shortURL.UpdatedAt = time.Now().Format(time.RFC3339)

	_, err = db.Collection().UpdateOne(ctx, bson.M{"shortCode": shortCode}, bson.M{"$set": shortURL})
	if err != nil {
		log.Printf("Failed to update short URL: %v", err)
		utils.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to update short URL", nil)
		return
	}

	utils.WriteJSON(c.Writer, http.StatusOK, "Short link updated successfully", shortURL)

}

// @Summary Delete a short URL
// @Description Delete a short URL by short code
// @Tags Shorten
// @Param shortCode path string true "Short Code"
// @Success 200 {string} string "Short link deleted successfully"
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/shorten/{shortCode} [delete]
func handleDeleteShortLink(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		utils.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid short code", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection().DeleteOne(ctx, bson.M{"shortCode": shortCode})
	if err != nil {
		log.Printf("Failed to delete short URL: %v", err)
		utils.WriteJSON(c.Writer, http.StatusInternalServerError, "Failed to delete short URL", nil)
		return
	}

	utils.WriteJSON(c.Writer, http.StatusOK, "Short link deleted successfully", nil)
}

// @Summary Get stats for a short URL
// @Description Retrieve statistics for a short URL
// @Tags Shorten
// @Param shortCode path string true "Short Code"
// @Success 200 {object} models.ShortURL
// @Failure 404 {object} models.Response
// @Router /api/v1/shorten/{shortCode}/stats [get]
func handleShortLinkStats(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		utils.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid short code", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shortURL models.ShortURL
	err := db.Collection().FindOne(ctx, bson.M{"shortCode": shortCode}).Decode(&shortURL)
	if err != nil {
		utils.WriteJSON(c.Writer, http.StatusNotFound, "Short URL not found", nil)
		return
	}

	utils.WriteJSON(c.Writer, http.StatusOK, "Short link stats", shortURL)
}
