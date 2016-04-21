package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func main() {
	// Setup
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)

}
