package webserver

import (
	"net/http"
	"time"

	"github.com/spotify-bot/server/internal/spotify"
	"github.com/labstack/echo/v4"
)

type IndexResponse struct {
	Version string `json:"version"`
}

func (s *WebServer) Index(c echo.Context) error {
	c.JSON(http.StatusOK, IndexResponse{Version: "v0.1.0"})
	return nil
}

func (s *WebServer) SpotifyConnect(c echo.Context) error {
	//Get Authorization code from spotify
	return nil
}

func (s *WebServer) SpotifyCallback(c echo.Context) error {
	//Receive access and refesh token from spotify
	platformCookie, err := c.Cookie("platform")
	if err != nil {
		s.server.Logger.Panic(err)
	}
	userCookie, err := c.Cookie("user_id")
	if err != nil {
		s.server.Logger.Panic(err)
	}

	var oathPlatform spotify.OauthPlatform
	switch platform := platformCookie.Value; platform {
	case "telegram":
		oathPlatform = spotify.PlatformTelegram
	default:
		s.server.Logger.Panic("Unsupported Platform")
	}

	userID := userCookie.Value
	authCode := c.QueryParam("code")
	s.spotify.AddUser(authCode, oathPlatform, userID) //TODO get error

	return c.String(200, "Authentication Successful")
}

func (s *WebServer) TelegramAuth(c echo.Context) error {

	userID := c.QueryParam("user_id")
	//TODO return error if user_id is not set

	// Set user_id cookie
	//TODO use JWT token to encrypt cookies
	userCookie := new(http.Cookie)
	userCookie.Name = "user_id"
	userCookie.Value = userID
	userCookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(userCookie)

	// Set Platform cookie
	platformCookie := new(http.Cookie)
	platformCookie.Name = "platform"
	platformCookie.Value = "telegram"
	platformCookie.Expires = time.Now().Add(1 * time.Hour)
	c.SetCookie(platformCookie)

	url := s.spotify.GetAuthURL()
	c.Redirect(http.StatusFound, url)
	return nil
}
