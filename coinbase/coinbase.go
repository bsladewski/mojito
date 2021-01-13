package coinbase

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bsladewski/mojito/candlestick"
	"github.com/bsladewski/mojito/data"
	"github.com/bsladewski/mojito/env"
	"github.com/sirupsen/logrus"
)

// init loads the configuration for connecting to the Coinbase API.
func init() {

	// retreive Coinbase configuration from environment
	websocketFeedBaseURL = env.GetStringSafe(websocketFeedBaseURLVariable,
		"wss://ws-feed.pro.coinbase.com")
	recordPrices = env.GetBoolSafe(recordPricesVariable, true)
	recordPricesInterval = env.GetIntSafe(recordPricesIntervalVariable, 60)

	// if we should record prices, spin up a process to listen for and aggregate
	// price data
	if recordPrices {
		if err := initializeRecordPrices(); err != nil {
			log.Fatalf("failed to connect to Coinbase websocket feed, Err: %s",
				err.Error())
		}
	}

}

const (
	// websocketFeedBaseURLVariable defines the environment variable for the URL
	// that will be used to establish a websocket connection to Coinbase feeds.
	websocketFeedBaseURLVariable = "MOJITO_COINBASE_WEBSOCKET_FEED_BASE_URL"
	// recordPricesVariable defines the environment variable for whether we
	// should record price data as candlesticks.
	recordPricesVariable = "MOJITO_COINBASE_RECORD_PRICES"
	// recordPricesIntervalVariable defines the environment variable that
	// determines how long we should aggregate price data before recording a
	// candlestick.
	recordPricesIntervalVariable = "MOJITO_COINBASE_RECORD_PRICES_INTERVAL"
)

// websocketFeedBaseURL the URL that will be used to establish a websocket
// connection to Coinbase feeds.
var websocketFeedBaseURL string

// recordPrices indicates that we should record price data as candlesticks.
var recordPrices bool

// recordPricesInterval determines how many seconds we should aggregate data
// before recording a candlestick.
var recordPricesInterval int

// priceData is used to read price data from the Coinbase ticker feed.
type priceData struct {
	Type      string    `json:"type"`
	TradeID   int       `json:"trade_id"`
	Sequence  int64     `json:"sequence"`
	Time      time.Time `json:"time"`
	ProductID string    `json:"product_id"`
	Price     string    `json:"price"`
	Side      string    `json:"side"`
	LastSize  string    `json:"last_size"`
	BestBid   string    `json:"best_bid"`
	BestAsk   string    `json:"best_ask"`
}

// initializeRecordPrices .
func initializeRecordPrices() error {

	_, channel, err := WebsocketSubscribe(FeedTicker)
	if err != nil {
		return err
	}

	go func() {

		intervalStart := time.Now()

		candlesticks := map[string]candlestick.Candlestick{}

		for {

			var price priceData

			// read price data from feed
			if err := json.Unmarshal(<-channel, &price); err != nil {
				logrus.Error(err)
			}

			// skip any messages that aren't for ticker data
			if price.Type != "ticker" {
				continue
			}

			// get the ticker from the price data
			ticker := strings.Split(price.ProductID, "-")[0]

			// parse price from price data
			currentPrice, err := strconv.ParseFloat(price.Price, 64)
			if err != nil {
				logrus.Errorf("%v: %v", err, price)
				continue
			}

			// ensure we have a candlestick to track the price data for this
			// ticker
			if _, ok := candlesticks[ticker]; !ok {
				candlesticks[ticker] = candlestick.Candlestick{
					CreatedAt: intervalStart,
					Ticker:    ticker,
				}
			}

			// aggregate price data into current interval
			candlesticks[ticker] = candlesticks[ticker].Add(
				0.0, 0.0, 0.0, 0.0, 1)

			if candlesticks[ticker].Open == 0.0 {
				candlesticks[ticker] = candlesticks[ticker].SetOpen(currentPrice)
			}

			if currentPrice > candlesticks[ticker].High {
				candlesticks[ticker] = candlesticks[ticker].SetHigh(currentPrice)
			}

			if candlesticks[ticker].Low == 0 || currentPrice < candlesticks[ticker].Low {
				candlesticks[ticker] = candlesticks[ticker].SetLow(currentPrice)
			}

			// check if we should record the current interval
			if time.Now().Sub(intervalStart) >=
				time.Duration(recordPricesInterval)*time.Second {

				// set the close price and save the ticker to the database
				candlesticks[ticker] = candlesticks[ticker].SetClose(currentPrice)
				if err := candlestick.Save(context.Background(), data.DB(),
					candlesticks[ticker]); err != nil {
					logrus.Error(err)
				}

				delete(candlesticks, ticker)
				intervalStart = intervalStart.Add(
					time.Duration(recordPricesInterval) * time.Second)
			}

		}

	}()

	return nil

}