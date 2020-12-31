package webserver

import (
	"net/http"

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
	return nil
}

func (s *WebServer) SpotifyCallback(c echo.Context) error {
	return nil
}
