package main

import "net/http"

func pong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}
