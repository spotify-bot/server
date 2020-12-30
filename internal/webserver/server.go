package webserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type WebServer struct {
	server *echo.Echo
}

func New() *WebServer {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", Index)

	return &WebServer{e}
}

func (w *WebServer) Start(address string) error {
	return w.server.Start(address)
}
