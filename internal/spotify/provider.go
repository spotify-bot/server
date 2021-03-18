package spotify

import (
	"context"
	"encoding/json"
	"github.com/spotify-bot/server/internal/mongo"
	"github.com/spotify-bot/server/pkg/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
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

func (s *SpotifyProvider) AddUser(code string, platform spotify.OauthPlatform, userID string) error {
	ctx := context.Background() //FIXME
	token, err := s.authConfig.Exchange(ctx, code)
	if err != nil {
		return err
	}

	mongoRow := mongo.OAuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.Type(),
		Expiry:       token.Expiry,
		Platform:     string(platform),
		UserID:       userID,
	}
	if err = s.db.UpsertOAuthToken(ctx, mongoRow); err != nil {
		log.Println("Failed to add token to database")
	}

	log.Printf("Code: [%s]\nPlatform: [%s]\nUser ID: [%s]", token.AccessToken, platform, userID)
	return nil
}

func (s *SpotifyProvider) getUserToken(platform spotify.OauthPlatform, userID string) (*oauth2.Token, error) {

	ctx := context.Background() //TODO add timeout
	mongoToken, err := s.db.GetOAuthTokenByUserID(ctx, userID, string(platform))
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  mongoToken.AccessToken,
		RefreshToken: mongoToken.RefreshToken,
		TokenType:    mongoToken.TokenType,
		Expiry:       mongoToken.Expiry,
	}

	// get token source and retreive token again to make sure token is refereshed if needed
	tokenSource := s.authConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func getCurrentlyPlayingSong(client *http.Client) (*spotify.Track, error) {
	resp, err := client.Get(spotify.CurrentlyPlayingEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, spotify.CallbackError{
			spotify.CurrentlyPlayingEndpoint,
			resp.StatusCode,
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response spotify.CurrentlyPlayingResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response.Track, nil
}
