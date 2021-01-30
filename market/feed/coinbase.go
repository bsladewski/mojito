package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"mojito/data"
	"mojito/market"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// exchangeCoinbase is the name that will be used as the exchange name for any
// candlesticks created through the Coinbase feed.
const exchangeCoinbase = "COINBASE"

// coinbaseFeed is used to stream price data from the Coinbase API.
type coinbaseFeed struct {
	mutex        *sync.Mutex
	conn         *websocket.Conn
	interval     time.Duration
	candlesticks map[string]market.Candlestick
	channels     map[string]chan market.Candlestick
	close        bool
}

// coinbasePriceData is used to read price data from the Coinbase ticker feed.
type coinbasePriceData struct {
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

func (c *coinbaseFeed) GetChannel(exchange,
	ticker string) (chan market.Candlestick, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := formatCoinbaseFeedKey(exchange, ticker)

	// retrieve the channel for this ticker
	channel, ok := c.channels[key]
	if !ok {
		// if the channel does not exist create it now
		channel = make(chan market.Candlestick)
		c.channels[key] = channel
	}

	return channel, nil
}

func (c *coinbaseFeed) Check(exchange,
	ticker string) (market.Candlestick, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// retrieve the candlestick for this ticker
	candlestick, ok := c.candlesticks[formatCoinbaseFeedKey(exchange, ticker)]
	if !ok {
		return market.Candlestick{}, ErrTickerNotFound
	}

	// check that the candlestick contains data
	if candlestick.Volume == 0 {
		return market.Candlestick{}, ErrNoPriceData
	}

	return candlestick, nil
}

func (c *coinbaseFeed) Commit(exchange,
	ticker string) (market.Candlestick, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// format the key that will be used to get the candlestick associated with
	// the specified ticker
	key := formatCoinbaseFeedKey(exchange, ticker)

	// retrieve the candlestick for this ticker
	candlestick, ok := c.candlesticks[key]
	if !ok {
		return market.Candlestick{}, ErrTickerNotFound
	}

	// check that the candlestick contains price data
	if candlestick.Volume == 0 {
		return market.Candlestick{}, ErrNoPriceData
	}

	// check if this candlestick opens a new hour or a new day
	last, err := market.GetLastByTicker(context.Background(), data.DB(),
		strings.ToUpper(exchange), strings.ToUpper(ticker))
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Error(err)
	} else {

		if last.CreatedAt.Hour() != candlestick.CreatedAt.Hour() {
			candlestick = candlestick.SetOpensHour(true)
		}

		if last.CreatedAt.Day() != candlestick.CreatedAt.Day() {
			candlestick = candlestick.SetOpensDay(true)
		}

	}

	// save the candlestick
	if err := market.SaveCandlestick(context.Background(), data.DB(),
		candlestick); err != nil {
		logrus.Error(err)
	} else {
		logrus.Debugf("new candlestick: %v", candlestick)
	}

	// clear the candlestick data associated with this ticker
	c.candlesticks[key] = market.Candlestick{
		CreatedAt: time.Now(),
		Exchange:  strings.ToUpper(exchange),
		Ticker:    strings.ToUpper(ticker),
	}

	// send the candlestick to the candlestick channel if it exists
	if channel, ok := c.channels[key]; ok {
		channel <- candlestick
	}

	// return the final candlestick data
	return candlestick, nil
}

func (c *coinbaseFeed) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.close = true
	return nil
}

// aggregate updates the current candlestick will additional price data. Returns
// whether the candlestick should be committed after aggregating the current
// price data.
func (c *coinbaseFeed) aggregate(exchange, ticker string,
	priceData coinbasePriceData) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// format the key that will be used to get the candlestick associated with
	// the specified ticker
	key := formatCoinbaseFeedKey(exchange, ticker)

	// parse price from price data
	currentPrice, err := strconv.ParseFloat(priceData.Price, 64)
	if err != nil {
		logrus.Errorf("%v: %v", err, priceData)
		return false
	}

	// retrieve the candlestick for this ticker
	candlestick, ok := c.candlesticks[key]
	if !ok {
		// if the candlestick is not found, initialize it now
		candlestick = market.Candlestick{
			CreatedAt: time.Now(),
			Exchange:  strings.ToUpper(exchangeCoinbase),
			Ticker:    strings.ToUpper(ticker),
		}
	}

	// increment the volume
	candlestick = candlestick.Add(0, 0, 0, 0, 1)

	// set the close price
	candlestick = candlestick.SetClose(currentPrice)

	if candlestick.Open == 0.00 {
		// if the open price has not been set, set it now
		candlestick = candlestick.SetOpen(currentPrice)
	}

	if candlestick.Low == 0.00 || candlestick.Low > currentPrice {
		// if the current price is lower than our low price or the low price
		// has not been set, set it now
		candlestick = candlestick.SetLow(currentPrice)
	}

	if candlestick.High == 0.00 || candlestick.High < currentPrice {
		// if the current price is higher than our high price or the high price
		// has not been set, set it now
		candlestick = candlestick.SetHigh(currentPrice)
	}

	// update the candlestick for this ticker
	c.candlesticks[key] = candlestick

	// check if we should commit this candlestick
	return time.Now().After(candlestick.CreatedAt.Add(c.interval))
}

// connectCoinbaseFeed connects to a feed of price data through the Coinbase
// API.
func connectCoinbaseFeed(platform feedPlatform) (Feed, error) {

	productIDList := []string{}

	// build a list of Coinbase product ids from the platform spec
	for _, security := range platform.Securities {
		productID := fmt.Sprintf("%s-%s",
			strings.ToUpper(security.Ticker),
			strings.ToUpper(security.ReferenceCurrency))
		productIDList = append(productIDList, productID)
	}

	// create the payload that will be used to initialize the connection to the
	// Coinbase ticker feed
	var subscribeMessage = struct {
		Type       string   `json:"type"`
		ProductIDs []string `json:"product_ids"`
		Channels   []string `json:"channels"`
	}{
		Type:       "subscribe",
		ProductIDs: productIDList,
		Channels:   []string{"ticker"},
	}

	// connect to the API websocket
	conn, _, err := websocket.DefaultDialer.Dial(platform.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	// send the subscribe message
	if err := conn.WriteJSON(subscribeMessage); err != nil {
		conn.Close()
		return nil, err
	}

	feed := &coinbaseFeed{
		mutex:        &sync.Mutex{},
		conn:         conn,
		interval:     platform.Interval,
		candlesticks: map[string]market.Candlestick{},
		channels:     map[string]chan market.Candlestick{},
	}

	// spawn a goroutine that continuously reads messages from the feed
	go func() {

		for !feed.close {

			// read a message from the feed
			_, message, err := conn.ReadMessage()
			if err != nil {
				logrus.Error(err)

				for {
					// if we run into an error attempting to read from the
					// feed, close the existing connection and attempt to
					// re-establish it after waiting for one minute
					feed.conn.Close()
					time.Sleep(time.Minute)
					if feed.close {
						break
					}

					feed.conn, _, err = websocket.DefaultDialer.Dial(
						platform.BaseURL, nil)
					if err != nil {
						logrus.Error(err)
					} else {
						break
					}
				}

				continue
			}

			var priceData coinbasePriceData

			// parse price data from the message
			if err := json.Unmarshal(message, &priceData); err != nil {
				logrus.Error(err)
			}

			// skip any messages that aren't for ticker data
			if priceData.Type != "ticker" {
				continue
			}

			// get the ticker from the price data
			ticker := strings.ToUpper(strings.Split(priceData.ProductID, "-")[0])

			// aggregate the price data and, if necessary, commit the current
			// candlestick
			if feed.aggregate(exchangeCoinbase, ticker, priceData) {
				_, err := feed.Commit(exchangeCoinbase, strings.ToUpper(ticker))
				if err != nil {
					logrus.Error(err)
				}
			}

		}

		conn.Close()

	}()

	return feed, nil
}

// formatCoinbaseFeedKey formats the supplied exchange and ticker into the
// format that is used to track different securities in the Coinbase feed.
func formatCoinbaseFeedKey(exchange, ticker string) string {
	return strings.ToUpper(fmt.Sprintf("%s-%s", exchange, ticker))
}
