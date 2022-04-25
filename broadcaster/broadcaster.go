package broadcaster

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/r3labs/sse/v2"
)

type Broadcaster struct {
	Server *sse.Server
}

type teamNameKey struct{}

var teamKey teamNameKey

func Serve() {
	server := sse.New()
	server.CreateStream("messages")
	b := Broadcaster{Server: server}
	b.serve()
}

func (b *Broadcaster) serve() {
	r := mux.NewRouter()

	subscribe := r.PathPrefix("/subscribe").Subrouter()
	subscribe.Use(b.validateClientJwt)
	subscribe.Use(b.loginAndSubscribe)
	subscribe.HandleFunc("", b.subscribeHandle)

	events := r.PathPrefix("/events").Subrouter()
	events.Use(b.validatePubSubJwt)
	events.HandleFunc("", b.pubSubHandle)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func (b *Broadcaster) subscribeHandle(w http.ResponseWriter, r *http.Request) {
	go func() {
		<-r.Context().Done()
		log.Printf("Client %s disconnected", r.RemoteAddr)
	}()

	b.Server.ServeHTTP(w, r)
}

func (b *Broadcaster) validateClientJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := strings.Split(r.Header.Get("Cookie"), ";")
		jwt := getJwtFromCookie(cookie)
		if jwt == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		claims, ok := validateAchilioJWT(jwt)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), teamKey, claims["hd"])
		r = r.WithContext(ctx)
		fmt.Println("teamName is", r.Context().Value(teamKey))
		next.ServeHTTP(w, r)
	})
}

func (b *Broadcaster) validatePubSubJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("JWT VALIDATION OF PUBSUB SERVICE ACCOUNT")
		next.ServeHTTP(w, r)
	})
}

// Login and set the dedicated stream
func (b *Broadcaster) loginAndSubscribe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		teamName := r.Context().Value(teamKey).(string)
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
	e := &SseEvent{Event: m.Message.Attributes.Type, ProjectId: m.Message.Attributes.ProjectID}
	data, err := json.Marshal(e)
	if err != nil {
		log.Printf("Error while serializing event response: %v", err)
	}
	b.Server.Publish(m.Message.Attributes.TeamName, &sse.Event{Data: data})
}
