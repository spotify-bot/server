package webserver

import (
	"github.com/spotify-bot/server/internal/spotify"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type spotifyProxy struct {
	reverseProxy *httputil.ReverseProxy
	spotify      *spotify.SpotifyProvider
}

func (s *WebServer) ProxyRequest() *spotifyProxy {

	target, err := url.Parse("https://api.spotify.com")
	if err != nil {
		log.Fatal(err)
	}
	rp := httputil.NewSingleHostReverseProxy(target)
	return &spotifyProxy{
		reverseProxy: rp,
		spotify:      s.spotify,
	}
}

func (s *spotifyProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//TODO important: path validations and error
	splittedPath := strings.SplitAfterN(req.URL.Path, "/", 5)
	platform := strings.Trim(splittedPath[2], "/")
	userid := strings.Trim(splittedPath[3], "/")
	spotifyApiPath := "/" + splittedPath[4]
	log.Println("sending request for ", spotifyApiPath)

	var oauthPlatform spotify.OauthPlatform
	switch platform {
	case "telegram":
		oauthPlatform = spotify.PlatformTelegram
	default:
		log.Println("Unsupported Platform") //TODO Error handling
	}

	/*
		target, err := url.Parse("https://api.spotify.com")
		if err != nil {
			log.Fatal(err)
		}
	*/

	s.spotify.SetRequestHeader(req, oauthPlatform, userid)

	url := "https://api.spotify.com/" + strings.SplitAfterN(req.URL.Path, "/", 5)[4]

	request, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		log.Fatal("Failed to Create New Request: ", err)
	}

	//req.URL.Scheme = target.Scheme
	//req.URL.Host = target.Host
	//req.URL.Path = spotifyApiPath

	s.spotify.SetRequestHeader(req, oauthPlatform, userid)
	s.reverseProxy.ServeHTTP(rw, request)
}
