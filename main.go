// dbs2go - An example code how to write data-base based app
//
// Copyright (c) 2016 - Valentin Kuznetsov <vkuznet@gmail.com>
//
package main

import (
	"flag"
	"github.com/vkuznet/gomil/web"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8989", "server port number")
	var base string
	flag.StringVar(&base, "base", "", "base url")
	flag.Parse()
	web.Server(base, port)
}
