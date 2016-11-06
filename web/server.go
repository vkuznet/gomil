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
	"time"

	"github.com/rcrowley/go-metrics"
)

// profiler, see https://golang.org/pkg/net/http/pprof/
import _ "net/http/pprof"

type Metrics struct {
	Meter metrics.Meter
}

// web server
func Server(base, port string) {
	log.Printf("Start server localhost:%s/%s", port, base)
	http.HandleFunc(fmt.Sprintf("/%s", base), TaskHandler)

	// register metrics
	r := metrics.DefaultRegistry
	m := metrics.GetOrRegisterMeter("requests", r)
	go metrics.Log(r, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	serverMetrics := Metrics{Meter: m}

	// start dispatcher for incoming requests
	var maxWorker, maxQueue int
	var err error
	maxWorker, err = strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		maxWorker = 10
	}
	maxQueue, err = strconv.Atoi(os.Getenv("MAX_QUEUE"))
	if err != nil {
		maxQueue = 100
	}
	dispatcher := NewDispatcher(maxWorker, maxQueue, serverMetrics)
	dispatcher.Run()
	log.Println("Start dispatcher with", maxWorker, "workers, queue size", maxQueue)

	// start server
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
