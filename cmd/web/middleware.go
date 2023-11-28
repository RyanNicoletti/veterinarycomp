package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
)

// add CSRF protection for post requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{HttpOnly: true, Path: "/", Secure: app.IsProduction, SameSite: http.SameSiteLaxMode})
	return csrfHandler
}

// load and save session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageStr := chi.URLParam(r, "page")
		var page int

		if pageStr != "" {
			page, _ = strconv.Atoi(pageStr)
		} else {
			page = 1
		}
		ctx := context.WithValue(r.Context(), "page", page)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
