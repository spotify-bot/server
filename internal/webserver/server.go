package webserver

import (
	"github.com/spotify-bot/server/internal/spotify"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.GET("/auth/connect", webServer.SpotifyConnect)   //Get Authorization code from spotify
	e.GET("/auth/callback", webServer.SpotifyCallback) //Receive access and refesh token from spotify
	e.GET("/auth/telegram", webServer.TelegramAuth)    //Authenticate Telegram user to spotify

	return webServer
}

func (w *WebServer) Start(address string) error {
	return w.server.Start(address)
}
