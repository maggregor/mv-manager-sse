package broadcaster

import (
	"encoding/json"
)

type MessageContent struct {
	Sa string `json:"serviceAccount"`
}

// UnmarshalJSON Custom unmarshall method for Attributes to check that the payload is valid for terraform executor
func (attribute *Attributes) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &attribute)
	if err != nil {
		return err
	}
	return nil
}
