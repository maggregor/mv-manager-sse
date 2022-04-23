package broadcaster

import (
	"encoding/json"
	"io/ioutil"
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

	subscribe := r.PathPrefix("/subscribe").Subrouter()
	subscribe.Use(b.preHandle1)
	subscribe.Use(b.loginAndSubscribe)
	subscribe.HandleFunc("", b.handle)

	events := r.PathPrefix("/events").Subrouter()
	events.Use(b.preHandle2)
	events.HandleFunc("", b.pubSubHandle)
	log.Fatal(http.ListenAndServe(":8080", r))
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
		log.Println("JWT VALIDATION OF CLIENT")
		next.ServeHTTP(w, r)
	})
}

func (b *Broadcaster) preHandle2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("JWT VALIDATION OF PUBSUB SERVICE ACCOUNT")
		next.ServeHTTP(w, r)
	})
}

// Login and set the dedicated stream
func (b *Broadcaster) loginAndSubscribe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		teamName := r.Header.Get("Cookie")
		log.Printf("Client %s logged succesfully", r.RemoteAddr)
		// Create dedicated stream if it doesn't exists
		if !b.Server.StreamExists(teamName) && teamName != "" {
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

// PubSub receives and processes a Pub/Sub push message.
func (b *Broadcaster) pubSubHandle(w http.ResponseWriter, r *http.Request) {
	var m PubSubMessage
	body, err := ioutil.ReadAll(r.Body)
	log.Printf("%s", string(body))
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Printf("Executing message %v", m.Message.ID)
	b.Server.Publish(m.Message.Attributes.TeamName, &sse.Event{Data: []byte(m.Message.Attributes.Type)})
}
