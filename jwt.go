package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSecret = []byte("CHANGE_TOKEN")

func newJwt(w http.ResponseWriter, r *http.Request) {
	//basic jwt token with an expiration flag
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expires": time.Now().Add(5 * time.Minute).Unix(),
	})
	tokenStr, err := token.SignedString(hmacSecret)
	if err != nil {
		println(err.Error())
		w.Write([]byte("Could not generate jwt"))
	} else {
		w.Write([]byte(tokenStr))
	}
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if the request is going to the /jwt endpoint we just forward it because there can't be a jwt token inside the request
		//to verify
		if r.URL.String() == "/jwt" {
			next.ServeHTTP(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		//Split because format is Bearer <Token>
		arg := strings.Split(auth, " ")
		if len(arg) != 2 || arg[0] != "Bearer" {
			w.Write([]byte("Not a valid Token"))
			return
		}
		//parse token
		token, err := jwt.Parse(arg[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing %s", t.Header["alg"])
			}
			return hmacSecret, nil
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			//check if it is still valid
			if int64(claims["expires"].(float64)) < time.Now().Unix() {
				w.Write([]byte("Token expired"))
				return
			} else {
				next.ServeHTTP(w, r)
			}
		} else {
			return
		}
	})

}
