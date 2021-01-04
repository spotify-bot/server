package webserver

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"time"

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
	platform := platformCookie.Value
	userID := userCookie.Value
	authCode := c.QueryParam("code")

	ctx := context.Background() //FIXME
	token, err := s.authConfig.Exchange(ctx, authCode)
	if err != nil {
		s.server.Logger.Fatal(err)
	}

	client := s.authConfig.Client(ctx, token)
	client.Get("...") //FIXME
	//TODO strore to DB

	return c.String(http.StatusOK, platform+"\n"+userID+"\n"+token.AccessToken)
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

	url := s.authConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(302, url)
	return nil
}
