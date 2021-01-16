package delivery

import (
	"net/http"
	"strings"
	"time"

	"mojito/cache"
	"mojito/candlestick"
	"mojito/data"
	"mojito/httperror"
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
	exchangeList, err := candlestick.ListExchanges(c, data.DB())
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
		tickerList, err := candlestick.ListTickers(c, data.DB(), exchange)
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
				Name: candlestick.GetTickerName(exchange, ticker),
			})
		}

		// add exchange data to response
		response.Exchanges = append(response.Exchanges, candlestickSpecExchange{
			ID:      exchange,
			Name:    candlestick.GetExchangeName(exchange),
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

	// TODO: add options for applying functions to candlestick data

	// TODO: add a way to specify date range and sample size

	// retrieve candlestick data
	candlesticks, err := candlestick.ListByTicker(c, data.DB(), exchange,
		ticker, time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
		return
	}

	// respond with candlesticks
	c.JSON(http.StatusOK, []listCandlestickResponse{
		{
			Function:     "List",
			Candlesticks: candlesticks,
		},
	})
}
