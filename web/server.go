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
	Meter        metrics.Meter
	WorkerMeters []metrics.Meter
}

var ServerMetrics Metrics

// web server
func Server(base, port string) {
	log.Printf("Start server localhost:%s/%s", port, base)
	http.HandleFunc(fmt.Sprintf("/%s", base), TaskHandler)

	// register metrics
	r := metrics.DefaultRegistry
	m := metrics.GetOrRegisterMeter("requests", r)
	go metrics.Log(r, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	// start dispatcher for incoming requests
	var workerMeters []metrics.Meter
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

	for i := 0; i < maxWorker; i++ {
		wm := metrics.GetOrRegisterMeter(fmt.Sprintf("worker_%d", i), r)
		workerMeters = append(workerMeters, wm)
	}
	ServerMetrics = Metrics{Meter: m, WorkerMeters: workerMeters}

	dispatcher := NewDispatcher(maxWorker, maxQueue)
	dispatcher.Run()
	log.Println("Start dispatcher with", maxWorker, "workers, queue size", maxQueue)

	// start server
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
