package broadcaster

import (
	"encoding/json"
)

type MessageContent struct {
	Sa string `json:"serviceAccount"`
}

func (attribute *Attributes) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &attribute)
	if err != nil {
		return err
	}
	return nil
}
