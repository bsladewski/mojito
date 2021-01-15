package delivery

import "github.com/bsladewski/mojito/candlestick"

// listCandlestickResponse is used to format responses from the list
// candlesticks endpoint.
type listCandlestickResponse struct {
	Function     string                    `json:"function"`
	Candlesticks []candlestick.Candlestick `json:"candlesticks"`
}
