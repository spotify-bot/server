package webserver

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spotify-bot/server/internal/spotify"
)

type IndexResponse struct {
	Version string `json:"version"`
}

func (s *WebServer) Index(c echo.Context) error {
	c.JSON(http.StatusOK, IndexResponse{Version: "v0.1.0"})
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

func (s *WebServer) ReverseProxy(c echo.Context) error {
	platform := c.Param("platform")
	userid := c.Param("userid")

	//TODO this switch case can be moved to spotify/types.go as a helper
	var oauthPlatform spotify.OauthPlatform
	switch platform {
	case "telegram":
		oauthPlatform = spotify.PlatformTelegram
	default:
		s.server.Logger.Error("Unsupported Platform")
	}

	spotifyApiPath := "/" + strings.SplitAfterN(c.Request().URL.Path, "/", 5)[4]
	log.Println("sending request for ", spotifyApiPath)

	target, err := url.Parse("https://api.spotify.com")
	if err != nil {
		log.Fatal(err)
	}

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = spotifyApiPath
		s.spotify.SetRequestHeader(req, oauthPlatform, userid)
		log.Println("New Request: ", req.URL)
	}
	reverseProxy := httputil.ReverseProxy{Director: director}
	reverseProxy.ServeHTTP(c.Response(), c.Request())
	log.Println("Arrriiivvvessss")
	return nil
}

func (s *WebServer) ProxyRequest(c echo.Context) error {

	platform := c.Param("platform")
	userid := c.Param("userid")

	//TODO this switch case can be moved to spotify/types.go as a helper
	var oauthPlatform spotify.OauthPlatform
	switch platform {
	case "telegram":
		oauthPlatform = spotify.PlatformTelegram
	default:
		s.server.Logger.Error("Unsupported Platform")
	}

	spotifyApiPath := strings.SplitAfterN(c.Request().URL.Path, "/", 4)[3]
	u, err := url.Parse("https://api.spotify.com/" + spotifyApiPath)
	if err != nil {
		log.Println("Failed to parse url")
	}
	req := c.Request()
	req.URL = u
	log.Print("kiiiirrrr ", req.RequestURI)

	_, err = s.spotify.ProxyRequest(oauthPlatform, userid, req)
	if err != nil {
		s.server.Logger.Error(err)
		return c.String(http.StatusForbidden, "User has not authenticated")
	}
	return nil
}
