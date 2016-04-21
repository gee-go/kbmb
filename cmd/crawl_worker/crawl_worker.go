package main

import (
	"net/http"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gee-go/kbmb/cfg"
	"github.com/gee-go/kbmb/crawl"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

func main() {
	// Setup
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)

	nsqConfig := &cfg.NSQConfig{
		NSQDHosts: []string{"localhost:4150"},
	}

	manager := crawl.NewManager(nsqConfig)

	worker := manager.NewWorker(8)
	defer worker.Stop()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Run(standard.New(":0"))
}
