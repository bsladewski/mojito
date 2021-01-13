// Package coinbase provides functions for interacting with the coinbase API.
//
// Environment:
//     MOJITO_COINBASE_WEBSOCKET_FEED_BASE_URL
//         string - the URL that will be used to establish connections to
//                  coinbase websocket feeds.
//                  Default: wss://ws-feed.pro.coinbase.com
//     MOJITO_COINBASE_RECORD_PRICES
//         bool - a flag that indicates whether we should record price data as
//                candlesticks.
//                Default: false
//     MOJITO_COINBASE_RECORD_PRICES_INTERVAL
//         int - determines how long in seconds we should aggregate price data
//               before recording a candlestick.
//               Default: 60
package coinbase
