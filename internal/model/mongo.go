package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"time"
)

const oauthTokenCollection = "oauth_tokens"

type MongoStorageOptions struct {
	DSN string
}

type MongoStorage struct {
	database *mongo.Database
}

func NewMongoStorage(ctx context.Context, opts MongoStorageOptions) (*MongoStorage, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(opts.DSN))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	cs, _ := connstring.ParseAndValidate(opts.DSN)
	database := client.Database(cs.Database)

	return &MongoStorage{
		database: database,
	}, nil
}

func (m *MongoStorage) UpsertOAuthToken(ctx context.Context, token OAuthToken) error {
	collection := m.database.Collection(oauthTokenCollection)

	filter := bson.D{
		{"platform", token.Platform},
		{"user_id", token.UserID},
	}
	update := bson.D{{"$set", token}}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStorage) GetOAuthTokenByUserID(ctx context.Context, userID string, platform string) (*OAuthToken, error) {
	collection := m.database.Collection(oauthTokenCollection)

	filter := bson.D{
		{"platform", platform},
		{"user_id", userID},
	}

	result := new(OAuthToken)
	err := collection.FindOne(ctx, filter).Decode(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
