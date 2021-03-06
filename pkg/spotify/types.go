package spotify

import "fmt"

type OauthPlatform string

const (
	OpenSpotifyTrackEndpoint               = "https://open.spotify.com/track/"
	TrackURIPrefix                         = "spotify:track:"
	CurrentlyPlayingEndpoint               = "/v1/me/player/currently-playing"
	RecentlyPlayedEndpoint                 = "/v1/me/player/recently-played"
	AddToQueueEndpoint                     = "/v1/me/player/queue"
	PlaySongEndpoint                       = "/v1/me/player/play"
	PlatformTelegram         OauthPlatform = "telegram"
)

type CurrentlyPlayingResponse struct {
	Track Track `json:"item"` // Spotify API returns a single item on currently playing but consists of track
}

type RecentlyPlayedResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	Track Track `json:"track"`
}

type Track struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Type  string `json:"type"`
	Album Album  `json:"album"`
}

type Album struct {
	Artists []Artist `json:"artists"`
	Images  []Image  `json:"images"`
}

type Artist struct {
	Name string `json:"name"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type CallbackError struct {
	Endpoint string
	Code     int
}

func (e CallbackError) Error() string {
	return fmt.Sprintf("API call to %v failed with %v error", e.Endpoint, e.Code)
}
