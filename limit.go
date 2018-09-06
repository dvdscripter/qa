package main

import (
	"net/http"
	"sync"
	"time"
)

type limit struct {
	max      int
	interval time.Duration
	counter  int
	mux      sync.Mutex
}

func (l *limit) Pass() bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.counter >= l.max {
		return false
	}
	l.counter++
	return true
}

func (l *limit) start() {
	for {
		select {
		case <-time.NewTicker(l.interval).C:
			l.mux.Lock()
			l.counter = 0
			l.mux.Unlock()
		}
	}
}

func (l *limit) toLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Pass() {
			http.Error(w, "Rate limiting", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
