package webserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type IndexResponse struct {
	Version string `json:"version"`
}

func Index(c echo.Context) error {
	c.JSON(http.StatusOK, IndexResponse{Version: "v0.1.0"})
	return nil
}
