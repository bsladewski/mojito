// Package health ...
package health

import (
	"net/http"

	"github.com/bsladewski/mojito/server"
	"github.com/gin-gonic/gin"
)

// init binds API endpoints for checking application health.
func init() {
	server.Router().GET(healthEndpoint, healthHandler)
}

const (
	// healthEndpoint the API endpoint that checks whether the server is able to
	// complete requests.
	healthEndpoint = "/health"
)

// healthHandler returns a 200 - OK response.
func healthHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}
