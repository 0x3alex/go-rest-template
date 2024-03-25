package main

import "net/http"

type middleware func(http.Handler) http.Handler

func middlewareStack(xs ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for _, v := range xs {
			next = v(next)
		}
		return next
	}
}
