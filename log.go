package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const apacheFormatPattern = "%v - - [%v] \"%v %v %v\" %.4f\n"

// ApacheLogRecord holds the relevant parts of an Apache-style log line
type ApacheLogRecord struct {
	http.ResponseWriter
	ip                    string
	time                  time.Time
	method, uri, protocol string
	status                int
	responseBytes         int64
	elapsedTime           time.Duration
}

// Log writes a line of Apache-style log to a Writer
func (r *ApacheLogRecord) Log(out io.Writer) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)
	_, err := fmt.Fprintf(out, apacheFormatPattern, r.ip, timeFormatted, requestLine, r.status, r.responseBytes,
		r.elapsedTime.Seconds())

	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func (r *ApacheLogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

// WriteHeader ...
func (r *ApacheLogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// ApacheLoggingHandler ...
type ApacheLoggingHandler struct {
	handler http.Handler
	out     io.Writer
}

// NewApacheLoggingHandler creates a new http.Handler
func NewApacheLoggingHandler(handler http.Handler, out io.Writer) http.Handler {
	return &ApacheLoggingHandler{
		handler: handler,
		out:     out,
	}
}

func (h *ApacheLoggingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	record := &ApacheLogRecord{
		ResponseWriter: rw,
		ip:             clientIP,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		elapsedTime:    time.Duration(0),
	}

	startTime := time.Now()
	h.handler.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	record.Log(h.out)
}
