package broadcaster

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/r3labs/sse/v2"
)

type Broadcaster struct {
	Server *sse.Server
}

func Serve() {
	server := sse.New()
	server.CreateStream("messages")
	b := Broadcaster{Server: server}
	b.serve()
}

func (b *Broadcaster) serve() {
	r := mux.NewRouter()

	r.HandleFunc("/", pubSubHandle)
	r.HandleFunc("/test", b.testHandler)
	r.HandleFunc("/events", b.handle)
	r.Use(b.preHandle1)
	r.Use(b.preHandle2)
	// r.Use(b.loginAndSubscribe)
	http.ListenAndServe(":8080", r)
}

func (b *Broadcaster) handle(w http.ResponseWriter, r *http.Request) {
	go func() {
		<-r.Context().Done()
		log.Printf("Client %s disconnected", r.RemoteAddr)
	}()

	b.Server.ServeHTTP(w, r)
}

func (b *Broadcaster) preHandle1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("STEP 1")
		b.Server.Publish("messages", &sse.Event{Data: []byte("ping")})
		next.ServeHTTP(w, r)
	})
}

func (b *Broadcaster) preHandle2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("STEP 2")
		next.ServeHTTP(w, r)
	})
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

func (b *Broadcaster) testHandler(w http.ResponseWriter, r *http.Request) {
	b.Server.Publish("messages", &sse.Event{Data: []byte("ping")})
}
