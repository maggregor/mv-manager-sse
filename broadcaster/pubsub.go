package broadcaster

// Attributes is the payload of the attributes field in the message of a Pub/Sub event.
type Attributes struct {
	TeamName  string `json:"teamName"`
	ProjectID string `json:"projectId"`
	Type      string `json:"eventType"`
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
