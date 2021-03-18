package webserver

import (
	"bytes"
	provider "github.com/spotify-bot/server/internal/spotify"
	"github.com/spotify-bot/server/pkg/spotify"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type spotifyProxy struct {
	reverseProxy *httputil.ReverseProxy
	spotify      *provider.SpotifyProvider
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
	splittedPath := strings.SplitAfterN(req.URL.RequestURI(), "/", 5)
	platform := strings.Trim(splittedPath[2], "/")
	userid := strings.Trim(splittedPath[3], "/")
	spotifyApiPath := splittedPath[4]

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

	url := "https://api.spotify.com/" + spotifyApiPath

	newRequest, err := http.NewRequest(req.Method, url, nil)
	if err != nil {
		log.Fatal("Failed to Create New Request: ", err)
	}

	var b bytes.Buffer
	b.ReadFrom(req.Body)
	//req.Body = ioutil.NopCloser(&b)
	newRequest.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))

	//req.URL.Scheme = target.Scheme
	//req.URL.Host = target.Host
	//req.URL.Path = spotifyApiPath

	s.spotify.SetRequestHeader(newRequest, oauthPlatform, userid)
	log.Println("sending request for ", newRequest.URL.RequestURI())
	log.Println("old request body: ", req.ContentLength)
	log.Println("request body: ", newRequest.ContentLength)
	s.reverseProxy.ServeHTTP(rw, newRequest)
}
