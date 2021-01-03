package webserver

import (
	"github.com/koskalak/mamal/internal/mongo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type WebServerOptions struct {
	Mongo               *mongo.MongoStorage
	SpotifyClientID     string
	SpotifyClientSecret string
}

type WebServer struct {
	server              *echo.Echo
	mongo               *mongo.MongoStorage
	spotifyClientID     string
	spotifyClientSecret string
}

func New(opts WebServerOptions) *WebServer {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	webServer := &WebServer{
		server:              e,
		mongo:               opts.Mongo,
		spotifyClientID:     opts.SpotifyClientID,
		spotifyClientSecret: opts.SpotifyClientSecret,
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
