package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"` // Data is omitted if nil or empty
}

type ShortURL struct {
	ID          primitive.ObjectID `bson:"_id"`
	URL         string             `bson:"url"`
	ShortCode   string             `bson:"shortCode"`
	CreatedAt   string             `bson:"createdAt"`
	UpdatedAt   string             `bson:"updatedAt"`
	AccessCount int                `bson:"accessCount"`
}
