package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Version holds the short hash of the current git commit
var Version string

// getLocalIP returns the non loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func ipHandler(w http.ResponseWriter, r *http.Request) error {
	_, err := fmt.Fprintf(w, getLocalIP()+"\n")
	return err
}

func timeHandler(w http.ResponseWriter, r *http.Request) error {
	_, err := fmt.Fprintf(w, time.Now().String()+"\n")
	return err
}

func hostnameHandler(w http.ResponseWriter, r *http.Request) error {
	name, err := os.Hostname()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, name+"\n")
	return err
}

func versionHandler(w http.ResponseWriter, r *http.Request) error {
	var err error
	// w.WriteHeader(http.StatusInternalServerError)
	_, err = fmt.Fprintf(w, Version+"\n")
	return err
}

func fibonacciHandler(w http.ResponseWriter, r *http.Request) error {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		return err
	}

	query, _ := url.ParseQuery(u.RawQuery)
	iterStr := query["n"][0]
	iter, err := strconv.ParseUint(iterStr, 10, 32)
	if err != nil {
		return err
	}

	if iter > 40 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, err = fmt.Fprintf(w, "Only values of n <= 40 are supported.\n")
		return err
	}
	_, err = fmt.Fprintf(w, strconv.FormatUint(uint64(fibonacci(uint(iter))), 10)+"\n")
	return err
}

func fibonacci(n uint) uint {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fibonacci(n-1) + fibonacci(n-2)
	}
}

func errorHandler(handler func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func logHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	logger := NewApacheLoggingHandler(http.HandlerFunc(handler), os.Stderr)
	return func(w http.ResponseWriter, r *http.Request) {
		logger.ServeHTTP(w, r)
	}
}

func main() {
	log.Printf("%v (%v) starting up...", os.Args[0], Version)

	tracer.Start(
		tracer.WithService("webby"),
		tracer.WithServiceVersion("2.0"),
	)

	http.HandleFunc("/", logHandler(errorHandler(ipHandler)))
	http.HandleFunc("/time", logHandler(errorHandler(timeHandler)))
	http.HandleFunc("/hostname", logHandler(errorHandler(hostnameHandler)))
	http.HandleFunc("/version", logHandler(errorHandler(versionHandler)))
	http.HandleFunc("/fibonacci", logHandler(errorHandler(fibonacciHandler)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
