// Package server exposes a server router than can be used to bind API endpoints
// and provides functions for managing the server.
package server

import (
	"fmt"

	"github.com/bsladewski/mojito/env"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// init initializes the application server router.
func init() {
	if router == nil {
		router = gin.Default()
	}
}

const (
	// PortVariable defines the environment variable for the server port.
	portVariable = "MOJITO_PORT"
	// TLSCertVariable defines the environment variable for the TLS certificate.
	// If set the server will be run using TLS encryption.
	tlsCertVariable = "MOJITO_CERT"
	// TLSKeyVariable defines the environment variable for the TLS key. If set
	// the server will run using TLS encryption.
	tlsKeyVariable = "MOJITO_KEY"
	// HTTPDefaultPort the default port when running the server without TLS
	// encryption and no explicit port.
	httpDefaultPort = 80
	// HTTPSDefaultPort the default port when running the server with TLS
	// encryption and no explicit port.
	httpsDefaultPort = 443
)

// router is used to bind API endpoints.
var router *gin.Engine

// Router retrieves the application server router which can be used to bind
// handler functions to API endpoints.
func Router() *gin.Engine {
	return router
}

// Run starts the application server. Returns when the server is terminated.
func Run() {

	cert, key := env.GetString(tlsCertVariable), env.GetString(tlsKeyVariable)

	// check if we should be running the server using TLS encryption
	if cert != "" || key != "" {

		// run the server using HTTPS
		port := env.GetIntSafe(portVariable, httpsDefaultPort)

		logrus.Infof("starting HTTPS server on port %d", port)
		logrus.Error(router.RunTLS(
			fmt.Sprintf(":%d", port),
			cert, key,
		))

	} else {

		// run the server using HTTP
		port := env.GetIntSafe(portVariable, httpDefaultPort)

		logrus.Infof("starting HTTP server on port %d", port)
		logrus.Error(router.Run(
			fmt.Sprintf(":%d", port),
		))

	}

}
