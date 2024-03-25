package main

import "net/http"

func main() {
	router := http.NewServeMux()
	stack := middlewareStack(
		logging,
		auth,
	)
	router.HandleFunc("GET /jwt", newJwt)
	router.HandleFunc("GET /pong", pong)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}
	server.ListenAndServe()
}
