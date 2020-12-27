// Package main is the entry point for the mojito server application.
package main

import (
	"github.com/bsladewski/mojito/server"

	_ "github.com/bsladewski/mojito/health"
	_ "github.com/bsladewski/mojito/user/delivery"
)

// main stands up a mojito server.
func main() {

	// run the API server
	server.Run()

}
