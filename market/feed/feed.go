package feed

import (
	"errors"
	"mojito/market"
	"sync"

	"github.com/sirupsen/logrus"
)

// ErrNoPriceData is returned when a request is made for price data but the
// current interval has no volume.
var ErrNoPriceData = errors.New("no price data for this interval")

// ErrTickerNotFound is returned when a request is made for price data but the
// specified exchange and ticker are not found in the feed.
var ErrTickerNotFound = errors.New("ticker not found")

// Feed encapsulates a stream of market data.
type Feed interface {
	// GetChannel retrieves a channel that will receive a candlestick every time
	// the feed commits price data.
	GetChannel(exchange, ticker string) (chan market.Candlestick, error)
	// AddSecurity will subscribe the feed to a new ticker and return a channel
	// that will receive a candlestick every time the feed commits price data.
	AddSecurity(exchange, ticker string) (chan market.Candlestick, error)
	// Check retrieves the candlestick that is currently being aggregated.
	Check(exchange, ticker string) (market.Candlestick, error)
	// Commit saves the candlestick that is currently being aggregated and
	// begins aggregating a new candlestick.
	Commit(exchange, ticker string) (market.Candlestick, error)
	// Close commits the current candlestick and stops listening for new price
	// data.
	Close() error
}

// feeds keeps track of all feeds of price data.
var feeds = map[string]Feed{}

// mutex is used to facilitate concurrent access to the map of feeds.
var mutex = &sync.Mutex{}

// Connect establishes a new connection to the specified platform, if an
// existing connection already exists it will be closed and replaced by the new
// connection. If no platform integration exists this function will log a
// warning and return nil for the feed and error.
func Connect(platform *platformFeed) (Feed, error) {

	mutex.Lock()
	defer mutex.Unlock()

	if feed, ok := feeds[platform.Name]; ok {
		// if a connection to this feed already exists, close the feed and
		// remove it from the map of feeds
		if err := feed.Close(); err != nil {
			logrus.Error(err)
			delete(feeds, platform.Name)
		}
	}

	var feed Feed
	var err error

	// connect to the appropriate feed based on the platform name
	switch platform.Platform.Key {
	case market.PlatformCoinbase:
		feed, err = connectCoinbaseFeed(*platform)
		if err != nil {
			return nil, err
		}
	default:
		logrus.Warnf("feed \"%s\" is not defined", platform.Name)
	}

	if feed != nil {
		// if we were able to initialize the feed, add it to the map of feeds
		feeds[platform.Name] = feed
	}

	return feed, nil

}

// ptrToBool gets a pointer to the supplied boolean value.
func ptrToBool(val bool) *bool {
	return &val
}
