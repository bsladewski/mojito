package delivery

import (
	"net/http"
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
	server.Router().POST(listCandlestickEndpoint, user.JWTAuthMiddleware(),
		cache.LocalCacheMiddleware(60*time.Second), listCandlestick)

}

const (
	// listCandlestickEndpoint the API endpoint used to retrieve candlestick
	// data.
	listCandlestickEndpoint = "/candlestick/exchange/:exchange/ticker/:ticker"
)

// listCandlestick retrieves candlestick data.
func listCandlestick(c *gin.Context) {

	// read path parameters
	exchange := c.Param("exchange")
	ticker := c.Param("ticker")

	// retrieve candlestick data
	candlesticks, err := candlestick.ListByTicker(c, data.DB(), exchange,
		ticker, time.Now(), time.Now())
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, httperror.ErrorResponse{
			ErrorMessage: httperror.InternalServerError,
		})
	}

	// respond with candlesticks
	c.JSON(http.StatusOK, []listCandlestickResponse{
		{
			Function:     "List",
			Candlesticks: candlesticks,
		},
	})
}
