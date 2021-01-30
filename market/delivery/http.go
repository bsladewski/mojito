package delivery

import (
	"net/http"
	"strings"
	"time"

	"mojito/cache"
	"mojito/data"
	"mojito/httperror"
	"mojito/market"
	"mojito/server"
	"mojito/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// init registers the candlestick API with the application router.
func init() {

	// bind private endpoints
	server.Router().GET(candlestickSpecEndpoint, user.JWTAuthMiddleware(),
		cache.LocalCacheMiddleware(60*time.Second), candlestickSpec)
	server.Router().GET(listCandlestickEndpoint, user.JWTAuthMiddleware(),
		cache.LocalCacheMiddleware(60*time.Second), listCandlestick)

}

const (
	// candlestickSpecEndpoint the API endpoint used to retrieve available
	// options for requesting candlestick data.
	candlestickSpecEndpoint = "/candlestick/spec"
	// listCandlestickEndpoint the API endpoint used to retrieve candlestick
	// data.
	listCandlestickEndpoint = "/candlestick/exchange/:exchange/ticker/:ticker"
)

// candlestickSpec retrieves available options for requesting candlestick data.
func candlestickSpec(c *gin.Context) {

	// retrieve available exchanges
	exchangeList, err := market.ListExchanges(c, data.DB())
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	response := candlestickSpecResponse{
		Exchanges: []candlestickSpecExchange{},
	}

	// retrieve available tickers for each exchange
	for _, exchange := range exchangeList {

		// retrieve tickers
		tickerList, err := market.ListTickers(c, data.DB(), exchange)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
				ErrorMessage: httperror.InternalServerError,
			})
			return
		}

		tickers := []candlestickSpecTicker{}

		// get display names for tickers
		for _, ticker := range tickerList {
			tickers = append(tickers, candlestickSpecTicker{
				ID:   ticker,
				Name: ticker,
			})
		}

		// add exchange data to response
		response.Exchanges = append(response.Exchanges, candlestickSpecExchange{
			ID:      exchange,
			Name:    exchange,
			Tickers: tickers,
		})

	}

	//response with spec
	c.JSON(http.StatusOK, response)
}

// listCandlestick retrieves candlestick data.
func listCandlestick(c *gin.Context) {

	// read path parameters
	exchange := strings.ToUpper(c.Param("exchange"))
	ticker := strings.ToUpper(c.Param("ticker"))

	// default date range to the current day
	end := time.Now()
	start := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0,
		end.Location())

	// TODO: get start and end date from GET params

	var hourly bool
	var daily bool

	if end.Sub(start) > time.Hour*24*60 {
		// if the date range is larger than two months only retrieve
		// candlesticks that open a new day
		daily = true
	} else if end.Sub(start) > time.Hour*24 {
		// if the date range is larger than a day only retrieve candlesticks
		// that open a new hour
		hourly = true
	}

	// retrieve candlestick data
	candlesticks, err := market.ListByTicker(c, data.DB(), exchange,
		ticker, hourly, daily, start, end)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// respond with candlesticks
	c.JSON(http.StatusOK, candlesticks)
}
