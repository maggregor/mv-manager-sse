package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/r3labs/sse/v2"

	p "github.com/achilio/mv-manager-sse/broadcaster"
)

type Broadcaster struct {
	*sse.Server
}

func main() {
	server := sse.New()
	b := &Broadcaster{server}
	r := mux.NewRouter()
	r.HandleFunc("/", p.PubSub)
	r.HandleFunc("/events", b.handle)
	r.Use(b.loginAndSubscribe)
	http.ListenAndServe(":8080", r)
}

func (b *Broadcaster) handle(w http.ResponseWriter, r *http.Request) {
	go func() {
		<-r.Context().Done()
		log.Printf("Client %s disconnected", r.RemoteAddr)
		return
	}()

	b.Server.ServeHTTP(w, r)
}

// Login and set the dedicated stream
func (b *Broadcaster) loginAndSubscribe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		teamName := r.Header.Get("Cookie")
		log.Printf("Client %s logged succesfully", r.RemoteAddr)
		// Create dedicated stream if it doesn't exists
		if !b.Server.StreamExists(teamName) {
			b.Server.CreateStream(teamName)
			log.Printf("Stream for the team %s created", teamName)
		}
		// Add stream parameter in the URL for the r3labs/sse lib
		values := r.URL.Query()
		values.Set("stream", teamName)
		r.URL.RawQuery = values.Encode()
		log.Printf("Subscribe %s for the team %s", teamName, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
