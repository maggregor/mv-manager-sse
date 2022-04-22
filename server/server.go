package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/r3labs/sse/v2"
)

type Broadcaster struct {
	*sse.Server
}

func main() {
	server := sse.New()
	b := &Broadcaster{server}
	r := mux.NewRouter()
	r.HandleFunc("/events", b.handle)
	r.Use(b.loggingMiddleware)
	http.ListenAndServe(":8080", r)
}

func (b *Broadcaster) handle(w http.ResponseWriter, r *http.Request) {
	go func() {
		<-r.Context().Done()
		// Received Browser Disconnection
		return
	}()

	b.Server.ServeHTTP(w, r)
}

func (b *Broadcaster) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		teamName := r.Header.Get("Cookie")
		values := r.URL.Query()
		values.Set("stream", teamName)
		r.URL.RawQuery = values.Encode()
		b.Server.CreateStream(teamName)
		next.ServeHTTP(w, r)
	})
}

func (b *Broadcaster) CreateTeamStream(teamName string) {
	if !b.Server.StreamExists(teamName) {
		b.Server.CreateStream(teamName)
	}
}
