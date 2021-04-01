package main

import (
	"github.com/noissefnoc/kafun"
	"os"
)

func main() {
	cli := &kafun.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
