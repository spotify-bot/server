package spotify

type OauthPlatform string

const (
	OpenSpotifyTrackEndpoint               = "https://open.spotify.com/track/"
	RecentlyPlayedEndpoint                 = "https://api.spotify.com/v1/me/player/currently-playing"
	PlatformTelegram         OauthPlatform = "telegram"
)

type Response struct {
	Item Item `json:"item"`
}

type Item struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Type string `json:"type"`
}
