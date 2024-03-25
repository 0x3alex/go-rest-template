package main

import (
	"net/http"
	"time"
)

func logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()
		h.ServeHTTP(w, r)
		println(r.Method, r.URL.Path, time.Since(n))
	})
}
