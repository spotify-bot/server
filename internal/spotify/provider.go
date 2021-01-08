package spotify

import (
	"context"
	"github.com/koskalak/mamal/internal/mongo"
	"golang.org/x/oauth2"
	"log"
)

type ProviderOptions struct {
	DatabaseDSN string
	AuthConfig  *oauth2.Config
}

type SpotifyProvider struct {
	db         *mongo.MongoStorage
	authConfig *oauth2.Config
}

var provider *SpotifyProvider

func New(ctx context.Context, opts ProviderOptions) (*SpotifyProvider, error) {
	if provider != nil {
		return provider, nil
	}

	mongoStorage, err := mongo.NewMongoStorage(ctx, mongo.MongoStorageOptions{
		DSN: opts.DatabaseDSN,
	})
	if err != nil {
		return nil, err
	}

	return &SpotifyProvider{
		db:         mongoStorage,
		authConfig: opts.AuthConfig,
	}, nil
}

func (s *SpotifyProvider) GetAuthURL() string {
	return s.authConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (s *SpotifyProvider) AddUser(code, platform, userID string) error {
	ctx := context.Background() //FIXME
	token, err := s.authConfig.Exchange(ctx, code)
	if err != nil {
		return err
	}

	client := s.authConfig.Client(ctx, token)
	client.Get("...") //FIXME

	mongoRow := mongo.OAuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Platform:     platform,
		UserID:       userID,
	}
	if err = s.db.UpsertOAuthToken(ctx, mongoRow); err != nil {
		log.Println("Failed to add token to database")
	}

	//TODO strore to DB
	log.Println("Authentication Successful!")
	log.Printf("Code: [%s]\nPlatform: [%s]\nUser ID: [%s]", code, platform, userID)
	return nil
}
