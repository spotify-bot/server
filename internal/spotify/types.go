package spotify

type OauthPlatform string

const (
	OpenSpotifyTrackEndpoint               = "https://open.spotify.com/track/"
	TrackURIPrefix                         = "spotify:track:"
	CurrentlyPlayingEndpoint               = "https://api.spotify.com/v1/me/player/currently-playing"
	RecentlyPlayedEndpoint                 = "https://api.spotify.com/v1/me/player/recently-played"
	AddToQueueEndpoint                     = "https://api.spotify.com/v1/me/player/queue"
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
