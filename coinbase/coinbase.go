package coinbase

import "github.com/bsladewski/mojito/env"

// init loads the configuration for connecting to the Coinbase API.
func init() {

	// retreive Coinbase configuration from environment
	websocketFeedBaseURL = env.GetStringSafe(websocketFeedBaseURLVariable,
		"wss://ws-feed.pro.coinbase.com")

}

const (
	// websocketFeedBaseURLVariable defines the environment variable for the URL
	// that will be used to establish a websocket connection to Coinbase feeds.
	websocketFeedBaseURLVariable = "MOJITO_COINBASE_WEBSOCKET_FEED_BASE_URL"
)

// websocketFeedBaseURL the URL that will be used to establish a websocket
// connection to Coinbase feeds.
var websocketFeedBaseURL string
