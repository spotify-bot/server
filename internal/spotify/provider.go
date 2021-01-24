package spotify

import (
	"context"
	"encoding/json"
	"github.com/koskalak/mamal/internal/mongo"
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

func (s *SpotifyProvider) AddUser(code string, platform OauthPlatform, userID string) error {
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

func (s *SpotifyProvider) GetRecentlyPlayed(platform OauthPlatform, userID string) (track *Track, err error) {
	ctx := context.Background()
	mongoToken, err := s.db.GetOAuthTokenByUserID(ctx, userID, string(platform))
	if err != nil {
		return
	}

	token := &oauth2.Token{
		AccessToken:  mongoToken.AccessToken,
		RefreshToken: mongoToken.RefreshToken,
		TokenType:    mongoToken.TokenType,
		Expiry:       mongoToken.Expiry,
	}
	client := s.authConfig.Client(ctx, token)
	track, err = getCurrentlyPlayingSong(client)
	if err != nil {
		log.Println("no song currently playing")
		track, err = getRecentlyPlayedSong(client)
	}
	return
}

func getCurrentlyPlayingSong(client *http.Client) (*Track, error) {
	resp, err := client.Get(CurrentlyPlayingEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response CurrentlyPlayingResponse
	if err = json.Unmarshal(body, &response); err != nil {
		log.Println("No song playing atm, fetching recently played")
		return nil, err
	}
	return &response.Track, nil
}

func getRecentlyPlayedSong(client *http.Client) (*Track, error) {

	resp, err := client.Get(RecentlyPlayedEndpoint + "?limit=1") //FIXME
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response RecentlyPlayedResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response.Items[0].Track, nil
}
