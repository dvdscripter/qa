package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
	"securecodewarrior.com/ddias/heapoverflow/jwt"
)

type limit struct {
	*rate.Limiter
}

func (l *limit) toLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Allow() {
			http.Error(w, "Rate limiting", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func middleJSONLogger(fn appHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		toEncode := map[string]interface{}{}
		w.Header().Set("Content-Type", "application/json")
		payload := jwt.DecodePayload(r)

		resp, err := fn(w, r)
		if err != nil {
			toEncode["error"] = err.Error()
			toEncode["result"] = nil
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("E: %s %s %s %+v %s\n", r.RemoteAddr, r.Method,
				r.URL.Path, err, payload.Email)
		} else {
			toEncode["result"] = resp
			log.Printf("C: %s %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path,
				payload.Email)
		}
		json.NewEncoder(w).Encode(toEncode)
	})
}

func (app *app) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.isPublic(r) {
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Missing JWT", http.StatusInternalServerError)
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		rawToken := header[len("Bearer "):]
		token := jwt.NewFromFile(jwt.Payload{}, app.jwtKeyFile)
		if err := token.Decode(rawToken); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := token.Check(); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
