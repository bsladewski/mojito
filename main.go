// Package main is the entry point for the mojito server application.
// Environment:
//     MOJITO_ENABLE_DEBUG_LOG
//         bool - a flag that indicates whether the application should emit
//                debug level logs.
package main

import (
	"github.com/bsladewski/mojito/env"
	"github.com/bsladewski/mojito/server"
	"github.com/sirupsen/logrus"

	_ "github.com/bsladewski/mojito/candlestick/delivery"
	_ "github.com/bsladewski/mojito/health"
	_ "github.com/bsladewski/mojito/user/delivery"
)

const (
	// enableDebugLogVariable defines the environment variable that when set to
	// true will cause the application to emit debug level logs.
	enableDebugLogVariable = "MOJITO_ENABLE_DEBUG_LOG"
)

// main stands up a mojito server.
func main() {

	if env.GetBoolSafe(enableDebugLogVariable, false) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// run the API server
	server.Run()

}
