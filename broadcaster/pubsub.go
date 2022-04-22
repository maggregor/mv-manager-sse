package broadcaster

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Attributes is the payload of the attributes field in the message of a Pub/Sub event.
// For the Terraform executor to work for Achilio's Terraview, it needs to follow this structure
type Attributes struct {
	TeamName  string `json:"teamName"`
	ProjectID string `json:"projectId"`
	Type      string `json:"type"`
}

// Message is the payload of the message field of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type Message struct {
	Data       []byte     `json:"data,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
	ID         string     `json:"messageId"`
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Message      Message `json:"message"`
	Subscription string  `json:"subscription"`
}

// PubSub receives and processes a Pub/Sub push message.
func PubSub(w http.ResponseWriter, r *http.Request) {
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
}
