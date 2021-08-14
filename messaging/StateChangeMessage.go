package messaging

import (
	"encoding/json"
)

// as this is for MQs, I would not think it would need a swagger notation
type StateChangeMessage struct {
	BaseMessageData

	Id              int64
	PrevId          int64
	DisplayName     string
	PreDisplayName  string

	Record          interface{}
}

func (scm *StateChangeMessage) GetBody() []byte {
	b, _ := json.Marshal(scm)
	return b
}
