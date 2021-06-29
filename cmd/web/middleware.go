package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Hit the page")
// 		next.ServeHTTP(w, r)
// 	})
// }

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	cookie := http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	csrfHandler.SetBaseCookie(cookie)
	fmt.Println(cookie)
	return csrfHandler
}

// SessionLoad loads and saves the session on every requrest
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}