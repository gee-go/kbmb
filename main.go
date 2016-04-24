package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	apexCli "github.com/apex/log/handlers/cli"
	"github.com/gee-go/kbmb/cfg"
	"github.com/gee-go/kbmb/crawl"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/spf13/cobra"
)

func main() {
	log.SetHandler(apexCli.Default)
	log.SetLevel(log.DebugLevel)

	rootCmd := &cobra.Command{Use: "kbmb"}

	// Setup config
	config := &cfg.Cfg{}
	rootCmd.PersistentFlags().StringVar(&config.Redis.URL, "redis", "redis://redis:6379", "redis url")
	rootCmd.PersistentFlags().StringSliceVar(&config.NSQDHosts, "nsqd", []string{"kbmb_nsqd_1:4150", "kbmb_nsqd_2:4150", "kbmb_nsqd_3:4150"}, "nsqd hosts")

	rootCmd.AddCommand(&cobra.Command{
		Use: "worker",
		Run: func(c *cobra.Command, args []string) {
			manager := crawl.NewManager(config)
			worker := manager.NewWorker(8)
			defer worker.Stop()

			e := echo.New()
			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, World!")
			})
			e.Run(standard.New(":9999"))
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use: "start",
		Run: func(c *cobra.Command, args []string) {
			manager := crawl.NewManager(config)
			crawl, err := crawl.NewCrawl("mit.edu")
			if err != nil {
				panic(err)
			}
			manager.EmailConsumer(crawl, func(email string) {
				fmt.Println(email)
			})

			if err := manager.Start(crawl); err != nil {
				panic(err)
			}

			if err := manager.Wait(crawl); err != nil {
				panic(err)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
