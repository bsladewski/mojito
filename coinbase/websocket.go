package coinbase

import (
	"fmt"
	"sync"
	"time"

	"mojito/currency"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
)

// webSocketClient is used to keep track of subscribers to websocket feeds. As
// messages are received the raw message payloads will be sent to each
// subscriber to the feed.
type websocketClient struct {
	mutex *sync.Mutex
	feeds map[Feed]*websocketFeed
}

// websocketFeed listens for messages from a Coinbase websocket feed.
type websocketFeed struct {
	mutex       *sync.Mutex
	subscribers map[string]chan []byte
	conn        *websocket.Conn
}

// client is used to keep track of subscribers and is initialized when the first
// subscription to a feed is created.
var client *websocketClient

// Feed defines a support websocket feed.
type Feed string

const (
	// FeedTicker is used to retrieve realtime price updates for one or more
	// cryptocurrencies
	FeedTicker Feed = "ticker"
)

// WebsocketSubscribe subscribes to the target websocket feed and returns an id
// for the subscription and a channel that can be used to receive messages from
// the feed.
func WebsocketSubscribe(targetFeed Feed) (string, chan []byte, error) {

	// ensure the client has been initialized
	if client == nil {
		client = &websocketClient{
			mutex: &sync.Mutex{},
			feeds: map[Feed]*websocketFeed{},
		}
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	// ensure the target feed has been initialized
	feed, ok := client.feeds[targetFeed]
	if !ok {
		feed = &websocketFeed{
			mutex:       &sync.Mutex{},
			subscribers: map[string]chan []byte{},
		}

		// create the payload that will be used to initialize the connection to
		// the feed
		var subscribeMessage = struct {
			Type       string   `json:"type"`
			ProductIDs []string `json:"product_ids"`
			Channels   []string `json:"channels"`
		}{
			Type: "subscribe",
			ProductIDs: []string{
				fmt.Sprintf("%s-%s", currency.BTC, currency.BaseCurrency),
			},
			Channels: []string{
				string(targetFeed),
			},
		}

		// connect to the API websocket
		conn, _, err := websocket.DefaultDialer.Dial(websocketFeedBaseURL, nil)
		if err != nil {
			return "", nil, err
		}

		// send the subscribe message
		if err := conn.WriteJSON(subscribeMessage); err != nil {
			conn.Close()
			return "", nil, err
		}

		feed.conn = conn

		// spawn a goroutine that continuously reads messages from the feed
		go func() {

			for {

				// read a message from the feed
				_, message, err := conn.ReadMessage()
				if err != nil {
					logrus.Error(err)

					for {
						// re-establish the websocket connection
						feed.conn.Close()
						time.Sleep(time.Minute)

						feed.conn, _, err = websocket.DefaultDialer.Dial(
							websocketFeedBaseURL, nil)
						if err != nil {
							logrus.Error(err)
						} else {
							break
						}
					}

					continue
				}

				feed.mutex.Lock()

				// send the message to all subscribers of the feed
				for _, subscriber := range feed.subscribers {
					subscriber <- message
				}

				feed.mutex.Unlock()

			}

		}()

		client.feeds[targetFeed] = feed
	}

	feed.mutex.Lock()
	defer feed.mutex.Unlock()

	// generate the subscription id
	id := uuid.NewV4().String()

	// create the channel for passing along messages from the feed
	channel := make(chan []byte)

	// record the subscriber
	feed.subscribers[id] = channel

	return id, channel, nil

}

// WebsocketUnsubscribe finds and removes any subscriptions with the supplied
// id. The subscription channel will be closed by this operation.
func WebsocketUnsubscribe(id string) {

	if client == nil {
		// if the client hasn't been initialized don't do anything
		return
	}

	// check each feed for the subscriber
	for _, feed := range client.feeds {

		feed.mutex.Lock()

		if channel, ok := feed.subscribers[id]; ok {
			// if the subscriber id is found in the feed close its channel and
			// remove the subscriber from the feed
			close(channel)
			delete(feed.subscribers, id)
		}

		feed.mutex.Unlock()

	}

}
