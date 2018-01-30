package config

import (
	"github.com/urfave/cli"
)

var (
	FluentdAddress string
)

func Init(c *cli.Context) {
	FluentdAddress = c.String("fluentd-address")
}
