package main

import (
	"flag"
	"github.com/huoxue1/qinglong-go/cmd"
	"github.com/huoxue1/qinglong-go/service/config"
)

var (
	address string
)

func init() {
	flag.StringVar(&address, "add", "0.0.0.0:5700", "the ql listen address!")
	flag.Parse()
	config.SetAddress(address)
}

func main() {
	cmd.Run()
}
