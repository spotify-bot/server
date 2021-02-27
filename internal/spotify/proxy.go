package spotify

import (
	"errors"
	"log"
	"net/http"
)

func (s *SpotifyProvider) ProxyRequest(platform OauthPlatform, userid string, request *http.Request) (*http.Response, error) {
	client, err := s.getUserClient(platform, userid)
	if err != nil {
		log.Println("errrrrrrrr ", err)
		return nil, errors.New("User not found in db")
	}
	return client.Do(request)
}
