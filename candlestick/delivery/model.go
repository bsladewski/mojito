package delivery

import "mojito/candlestick"

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

// listCandlestickResponse is used to format responses from the list
// candlesticks endpoint.
type listCandlestickResponse struct {
	Function     string                    `json:"function"`
	Candlesticks []candlestick.Candlestick `json:"candlesticks"`
}
