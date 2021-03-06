package spotify

import (
	"errors"
	"log"
	"net/http"
)

// retired after reverse proxy
func (s *SpotifyProvider) ProxyRequest(platform OauthPlatform, userid string, request *http.Request) (*http.Response, error) {
	client, err := s.getUserClient(platform, userid)
	if err != nil {
		log.Println("errrrrrrrr ", err)
		return nil, errors.New("User not found in db")
	}
	return client.Do(request)
}

func (s *SpotifyProvider) SetRequestHeader(req *http.Request, platform OauthPlatform, userid string) {

	//FIXME what if user does not exist ? do not fill header so client gets the error from spotify ?
	token, err := s.getUserToken(platform, userid)
	if err != nil {
		return
	}
	token.SetAuthHeader(req)
}
