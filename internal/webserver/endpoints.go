package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/koskalak/mamal/config"
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
	for _, cookie := range c.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
	}
	authCode := c.QueryParam("code")
	return c.String(http.StatusOK, "code: "+authCode)
}

func (s *WebServer) TelegramAuth(c echo.Context) error {

	userID := c.QueryParam("user_id")
	//TODO return error if user_id is not set

	// Set user_id cookie
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

	//Authenticate Telegram user to spotify
	client_id := s.spotifyClientID
	redirect_uri := "http://" + config.AppConfig.Webserver.Address + "/auth/callback" //FIXME
	scope := "user-read-playback-state"
	url := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&scope=%s&redirect_uri=%s", client_id, scope, redirect_uri)

	// Redirect to Spotify auth URL
	c.Redirect(302, url)
	return nil
}
