// web server
// Copyright (c) 2015-2016 - Valentin Kuznetsov <vkuznet@gmail.com>
//
package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// profiler, see https://golang.org/pkg/net/http/pprof/
import _ "net/http/pprof"

// web server
func Server(base, port string) {
	log.Printf("Start server localhost:%s/%s", port, base)
	http.HandleFunc(fmt.Sprintf("/%s", base), TaskHandler)

	// start dispatcher for incoming requests
	maxWorker, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		maxWorker = 10
	}
	maxQueue, _ := strconv.Atoi(os.Getenv("MAX_QUEUE"))
	dispatcher := NewDispatcher(maxWorker, maxQueue)
	dispatcher.Run()
	log.Println("Start dispatcher with", maxWorker, "workers, queue size", maxQueue)

	// start server
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
