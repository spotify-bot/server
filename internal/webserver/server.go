package webserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spotify-bot/server/internal/spotify"

	"net/http/httputil"
)

type WebServerOptions struct {
	Spotify *spotify.SpotifyProvider
}

type WebServer struct {
	server  *echo.Echo
	spotify *spotify.SpotifyProvider
}

func New(opts WebServerOptions) *WebServer {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	webServer := &WebServer{
		server:  e,
		spotify: opts.Spotify,
	}

	e.GET("/", webServer.Index)
	e.GET("/auth/callback", webServer.SpotifyCallback) //Receive access and refesh token from spotify
	e.GET("/auth/telegram", webServer.TelegramAuth)    //Authenticate Telegram user to spotify

	// Spotify Endpoints
	e.Any("/spotify/:platform/:userid/*", echo.WrapHandler(webServer.ReverseProxy()))
	//TODO /spotify returns 404

	return webServer
}

func (w *WebServer) Start(address string) error {
	return w.server.Start(address)
}
