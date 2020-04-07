package api

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

type Error struct {
	ErrorCode string `json:"error_code"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

type StatusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *StatusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *StatusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Content-Type", "application/json")

		sw := StatusWriter{ResponseWriter: w}
		start := time.Now()

		// Necessary to pass cookies
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			return
		}

		next.ServeHTTP(&sw, r)
		duration := time.Now().Sub(start)

		// Round to milliseconds with 3 decimal places
		durationRounded := math.Round(duration.Seconds()*1000000) / 1000
		log.Println(fmt.Sprintf("status=%d method=%v url=%v duration=%v", sw.status, r.Method, r.RequestURI, durationRounded))
	})
}
