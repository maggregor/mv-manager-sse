package broadcaster

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/r3labs/sse/v2"
	"google.golang.org/api/idtoken"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type ConfigMap struct {
	JwtSecret string
	SAEmail   string
}

type Broadcaster struct {
	Server *sse.Server
	Config *ConfigMap
}

type teamNameKey struct{}

var teamKey teamNameKey

func Serve() {
	log.Println("starting server")
	server := sse.New()
	server.CreateStream("messages")
	c := config()
	b := Broadcaster{Server: server, Config: c}
	b.serve()
}

func config() *ConfigMap {
	var c ConfigMap
	c.JwtSecret = os.Getenv("JWT_SECRET")
	c.SAEmail = os.Getenv("SA_EMAIL")
	return &c
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
			return
		}
		claims, ok := validateAchilioJWT(jwt, b.Config.JwtSecret)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
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
		if r.Method != "POST" {
			log.Println("error invalid method")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		// Get the Cloud Pub/Sub-generated JWT in the "Authorization" header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(strings.Split(authHeader, " ")) != 2 {
			log.Println("error missing auth header")
			http.Error(w, "Missing Authorization header", http.StatusBadRequest)
			return
		}
		token := strings.Split(authHeader, " ")[1]
		log.Println(token)

		// Verify and decode the JWT.
		// If you don't need to control the HTTP client used you can use the
		// convenience method idtoken.Validate instead of creating a Validator.
		payload, err := idtoken.Validate(r.Context(), token, "https://dev.s.achilio.com")
		if err != nil {
			e := fmt.Sprintf("Invalid Token: %v", err)
			log.Println(e)
			log.Println(payload.Audience)
			http.Error(w, e, http.StatusBadRequest)
			return
		}
		log.Println(payload.Claims)
		if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
			log.Println("wrong issuer")
			http.Error(w, "Wrong Issuer", http.StatusUnauthorized)
			return
		}
		// IMPORTANT: you should validate claim details not covered by signature
		// and audience verification above, including:
		//   - Ensure that `payload.Claims["email"]` is equal to the expected service
		//     account set up in the push subscription settings.
		//   - Ensure that `payload.Claims["email_verified"]` is set to true.
		if payload.Claims["email"] != b.Config.SAEmail || payload.Claims["email_verified"] != true {
			log.Printf("unexpected email identity: got %v, expected %v", payload.Claims["email"], b.Config.SAEmail)
			http.Error(w, "Unexpected email identity", http.StatusUnauthorized)
			return
		}

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
