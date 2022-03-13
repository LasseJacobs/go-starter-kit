package middleware

import (
	"bufio"
	"errors"
	"github.com/rcrowley/go-metrics"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

// ResponseLatency returns a metric handler.
func ResponseLatency(next http.Handler) http.Handler {
	return CustomResponseLatency(metrics.DefaultRegistry)(next)
}

func CustomResponseLatency(r metrics.Registry) func(next http.Handler) http.Handler {
	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	err := r.Register("response_latency", h)
	if err != nil {
		logrus.Fatal(err)
	}

	okCount := metrics.NewCounter()
	err = r.Register("2XX_status", okCount)
	if err != nil {
		logrus.Fatal(err)
	}
	userErrCount := metrics.NewCounter()
	err = r.Register("4XX_status", userErrCount)
	if err != nil {
		logrus.Fatal(err)
	}
	sysErrCount := metrics.NewCounter()
	err = r.Register("5XX_status", sysErrCount)
	if err != nil {
		logrus.Fatal(err)
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			defer func() {
				h.Update(time.Since(t1).Nanoseconds())
			}()
			lrw := NewResponseWriter(w)
			next.ServeHTTP(lrw, r)

			if lrw.Status() >= 500 {
				sysErrCount.Count()
			} else if lrw.Status() >= 400 {
				userErrCount.Count()
			} else {
				okCount.Count()
			}
		}
		return http.HandlerFunc(fn)
	}
}

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 0 if the response has
	// not been written
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
	// Before allows for a function to be called before the ResponseWriter has been written to. This is
	// useful for setting headers or any other operations that must happen before a response has been written.
	Before(func(ResponseWriter))
}

type beforeFunc func(ResponseWriter)

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	nrw := &responseWriter{
		ResponseWriter: rw,
	}

	if _, ok := rw.(http.CloseNotifier); ok {
		return &responseWriterCloseNotifer{nrw}
	}

	return nrw
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []beforeFunc
}

func (rw *responseWriter) WriteHeader(s int) {
	if rw.Written() {
		return
	}
	rw.status = s
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Before(before func(ResponseWriter)) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		if !rw.Written() {
			// The status will be StatusOK if WriteHeader has not been called yet
			rw.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}

// Deprecated: the CloseNotifier interface predates Go's context package.
// New code should use Request.Context instead.
//
// We still implement it for backwards compatibliity with older versions of Go
type responseWriterCloseNotifer struct {
	*responseWriter
}

func (rw *responseWriterCloseNotifer) CloseNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
