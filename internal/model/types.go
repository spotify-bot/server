package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OAuthToken struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	AccessToken  string             `bson:"access_token"`
	RefreshToken string             `bson:"refresh_token"`
	TokenType    string             `bson:"token_type"`
	Expiry       time.Time          `bson:"expiry"`
	UserID       string             `bson:"user_id"`
	Platform     string             `bson:"platform"`
}

type OauthPlatform string

const (
	OauthPlatfomTelegram OauthPlatform = "telegram"
)
