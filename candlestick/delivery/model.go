package delivery

// candlestickSpecResponse is used to format responses from the get candlestick
// spec endpoint.
type candlestickSpecResponse struct {
	Exchanges []candlestickSpecExchange `json:"exchanges"`
}

// candlestickSpecExchange stores information about an available exchange.
type candlestickSpecExchange struct {
	ID      string                  `json:"id"`
	Name    string                  `json:"name"`
	Tickers []candlestickSpecTicker `json:"tickers"`
}

// candlestickSpecTicker stores information about an available ticker.
type candlestickSpecTicker struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
