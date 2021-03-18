package spotify

import (
	"github.com/spotify-bot/server/pkg/spotify"
	"net/http"
)

func (s *SpotifyProvider) SetRequestHeader(req *http.Request, platform spotify.OauthPlatform, userid string) {

	//FIXME what if user does not exist ? do not fill header so client gets the error from spotify ?
	token, err := s.getUserToken(platform, userid)
	if err != nil {
		return
	}
	token.SetAuthHeader(req)
}
