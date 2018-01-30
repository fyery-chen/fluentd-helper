package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/fluentd-helper/watcher"
	"github.com/urfave/cli"

	"github.com/rancher/fluentd-helper/config"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "fluentd-helper"
	app.Version = VERSION
	app.Usage = "fluentd helper use for monitory the fluentd config file change, and reload the fluentd use the latest configuration!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "fluentd-address",
			Usage: "Enable debug logs",
			Value: "127.0.0.1:24444",
		},
		cli.StringSliceFlag{
			Name:  "watched-file-list",
			Usage: "config file path list for what file need to be watched",
		},
	}
	app.Action = run

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	config.Init(c)

	wg := sync.WaitGroup{}
	jobs := make(chan int, 5)

	filePathList := c.StringSlice("watched-file-list")
	for _, v := range filePathList {
		wg.Add(1)
		go func(file string, jobs <-chan int) {
			defer wg.Done()
			logrus.Info("watched file:", file)
			watcher.Watcherfile(file, jobs)
		}(v, jobs)
	}

	waitForSignal()
	close(jobs)
	wg.Wait()

	fmt.Println("exiting")
	return nil
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	fmt.Println("receive exiting signal")

}
