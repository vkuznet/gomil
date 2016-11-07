package web

import (
	"log"
	"time"
)

// Processor is an object who process' given task
// The logic of the Processor should be implemented.
type Processor struct {
}

// Define execution process for a given task
func (e *Processor) Process(t *Task) error {
	return nil
}

var DefaultProcessor = &Processor{}

// A Request interface defines a task process
type Request interface {
	Process(*Task) error
}

// RequestFunc is a function type that implements the Request interface
type RequestFunc func(*Task) error

func (f RequestFunc) Process(t *Task) error {
	return f(t)
}

// A Decorator wraps a request with extra behavior
type Decorator func(Request) Request

// Logging returns a Decorator that logs client requests
func Logging(l *log.Logger) Decorator {
	return func(r Request) Request {
		return RequestFunc(func(t *Task) error {
			l.Println("Task", t)
			return r.Process(t)
		})
	}
}

// Pause returns a Decorator that pauses request for a given time interval
func Pause(interval time.Duration) Decorator {
	return func(r Request) Request {
		return RequestFunc(func(t *Task) error {
			log.Println("Task", t, "is paused by", interval)
			time.Sleep(interval)
			return r.Process(t)
		})
	}
}

// Tracer returns a Decorator that traces given request
func Tracer() Decorator {
	return func(r Request) Request {
		return RequestFunc(func(t *Task) error {
			log.Println("Trace", t)
			return r.Process(t)
		})
	}
}

// Decorate decorates a Request r with all given Decorators
func Decorate(r Request, ds ...Decorator) Request {
	decorated := r
	for _, decorate := range ds {
		decorated = decorate(decorated)
	}
	return decorated
}
